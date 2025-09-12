package main

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/go-ldap/ldap/v3"
)

// ================================
// 🔧 Конфигурация пула
// ================================
type LdapPoolConfig struct {
	URL         string        // ldap:// или ldaps://
	BindDN      string        // учетная запись для bind
	BindPass    string        // пароль
	MaxConns    int           // максимум соединений
	IdleTTL     time.Duration // через сколько закрывать idle-соединения
	DialTimeout time.Duration // таймаут для dial
	RetryCount  int           // количество попыток подключения при ошибке
	RetryDelay  time.Duration // пауза между попытками
}

// ================================
// 💡 Основная структура пула
// ================================
type LdapPool struct {
	config LdapPoolConfig
	conns  chan *pooledConn
	mu     sync.Mutex
	active int // сколько соединений создано
	stopGC chan struct{}
}

// 👇 обернутое соединение с временной меткой
type pooledConn struct {
	conn     *ldap.Conn
	lastUsed time.Time
}

// ================================
// ✅ Создание пула
// ================================
func NewLdapPool(cfg LdapPoolConfig) *LdapPool {
	pool := &LdapPool{
		config: cfg,
		conns:  make(chan *pooledConn, cfg.MaxConns),
		stopGC: make(chan struct{}),
	}

	// Запускаем очистку idle-соединений
	go pool.gcIdleConns()

	return pool
}

// ================================
// 📥 Получение соединения с timeout
// ================================
func (p *LdapPool) Get(timeout time.Duration) (*ldap.Conn, error) {
	// 1. Пытаемся взять из пула
	select {
	case pooled := <-p.conns:
		if isAlive(pooled.conn) {
			return pooled.conn, nil
		}
		// битое соединение — закрываем и создаем новое
		_ = pooled.conn.Close()
		p.decrActive()
		return p.newConnWithRetry()
	default:
		// 2. Если нет свободных, создаём новое (если можем)
		p.mu.Lock()
		if p.active < p.config.MaxConns {
			p.active++
			p.mu.Unlock()
			return p.newConnWithRetry()
		}
		p.mu.Unlock()

		// 3. Ждём освобождения соединения
		select {
		case pooled := <-p.conns:
			if isAlive(pooled.conn) {
				return pooled.conn, nil
			}
			_ = pooled.conn.Close()
			p.decrActive()
			return p.newConnWithRetry()
		case <-time.After(timeout):
			return nil, errors.New("LDAP pool: timeout waiting for connection")
		}
	}
}

// ================================
// ♻️ Вернуть соединение в пул
// ================================
func (p *LdapPool) Put(conn *ldap.Conn) {
	if conn == nil {
		return
	}
	select {
	case p.conns <- &pooledConn{conn: conn, lastUsed: time.Now()}:
	default:
		// если пул полон — закрываем соединение
		_ = conn.Close()
		p.decrActive()
	}
}

// ================================
// ❌ Закрыть пул
// ================================
func (p *LdapPool) Close() {
	close(p.stopGC)
	close(p.conns)
	for pooled := range p.conns {
		_ = pooled.conn.Close()
	}
}

// ================================
// 🔁 Фоновая очистка idle-соединений
// ================================
func (p *LdapPool) gcIdleConns() {
	ticker := time.NewTicker(p.config.IdleTTL / 2)
	defer ticker.Stop()

	for {
		select {
		case <-p.stopGC:
			return
		case <-ticker.C:
			p.cleanupIdle()
		}
	}
}

func (p *LdapPool) cleanupIdle() {
	now := time.Now()
	for {
		select {
		case pooled := <-p.conns:
			if now.Sub(pooled.lastUsed) > p.config.IdleTTL {
				_ = pooled.conn.Close()
				p.decrActive()
			} else {
				// вернуть обратно — оно еще живое
				select {
				case p.conns <- pooled:
				default:
					_ = pooled.conn.Close()
					p.decrActive()
				}
				return
			}
		default:
			return
		}
	}
}

// ================================
// 🔄 Подключение с Retry
// ================================
func (p *LdapPool) newConnWithRetry() (*ldap.Conn, error) {
	var conn *ldap.Conn
	var err error
	for i := 0; i < p.config.RetryCount; i++ {
		conn, err = ldap.DialURL(p.config.URL)
		if err == nil {
			err = conn.Bind(p.config.BindDN, p.config.BindPass)
			if err == nil {
				return conn, nil
			}
			conn.Close()
		}
		log.Printf("[ldap-pool] retry %d: %v", i+1, err)
		time.Sleep(p.config.RetryDelay)
	}
	p.decrActive()
	return nil, err
}

func (p *LdapPool) decrActive() {
	p.mu.Lock()
	p.active--
	p.mu.Unlock()
}

// ================================
// 🧪 Проверка соединения
// ================================
func isAlive(conn *ldap.Conn) bool {
	conn.SetTimeout(2 * time.Second)
	req := ldap.NewSearchRequest(
		"",
		ldap.ScopeBaseObject,
		ldap.NeverDerefAliases,
		1, 0, false,
		"(objectClass=*)",
		[]string{"dn"},
		nil,
	)
	_, err := conn.Search(req)
	return err == nil
}
