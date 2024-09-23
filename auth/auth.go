package auth

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/Cameron-Reed1/todo-web/types"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/scrypt"
)


var algorithm argon2idHasher = argon2idHasher{
    hashLen: 64,
    saltLen: 32,
    time: 6,
    memory: 24 * 1024,
    threads: 1,
}


func Hash(password, salt []byte) (*HashSalt, error) {
    return algorithm.Hash(password, salt)
}

func Validate(hash, salt, password []byte) bool {
    return algorithm.Validate(hash, salt, password)
}

func CreateSessionFor(user_id int64) (*types.Session, error) {
    buf := make([]byte, 32)
    _, err := rand.Read(buf)
    if err != nil {
        return nil, err
    }

    return &types.Session{ SessionId: base64.StdEncoding.EncodeToString(buf), UserId: user_id }, nil
}


func generateSalt(length uint) ([]byte, error) {
    salt := make([]byte, length)

    _, err := rand.Read(salt)
    if err != nil {
        return nil, err
    }

    return salt, nil
}

type HashSalt struct {
    Hash []byte
    Salt []byte
}

type hashAlgo interface {
    Hash(password, salt []byte) ([]byte, error)
    Validate(hash, salt, password []byte) bool
}

type scryptHasher struct {
    hashLen int
    saltLen uint
    cost int
    blockSize int
    parallelism int
}

type argon2idHasher struct {
    hashLen uint32
    saltLen uint
    time uint32
    memory uint32
    threads uint8
}

func (s *scryptHasher) Hash(password, salt []byte) (*HashSalt, error) {
    var err error

    if salt == nil || len(salt) == 0 {
        salt, err = generateSalt(s.saltLen)
        if err != nil {
            fmt.Println("\x1b[31mError: Failed to generate a password salt\x1b[0m")
            return nil, err
        }
    }

    hash, err := scrypt.Key(password, salt, s.cost, s.blockSize, s.parallelism, s.hashLen)
    if err != nil {
        return nil, err
    }

    return &HashSalt{Hash: hash, Salt: salt}, nil
}

func (s *scryptHasher) Validate(hash, salt, password []byte) bool {
    hashed_password, err := s.Hash(password, salt)
    if err != nil {
        fmt.Println("\x1b[31mError: Failed to generate a password hash\x1b[0m")
        return false
    }

    return bytes.Equal(hash, hashed_password.Hash)
}

func (a *argon2idHasher) Hash(password, salt []byte) (*HashSalt, error) {
    var err error

    if salt == nil || len(salt) == 0 {
        salt, err = generateSalt(a.saltLen)
        if err != nil {
            fmt.Println("\x1b[31mError: Failed to generate a password salt\x1b[0m")
            return nil, err
        }
    }

    hash := argon2.IDKey(password, salt, a.time, a.memory, a.threads, a.hashLen)

    return &HashSalt{Hash: hash, Salt: salt}, nil
}

func (s *argon2idHasher) Validate(hash, salt, password []byte) bool {
    hashed_password, err := s.Hash(password, salt)
    if err != nil {
        fmt.Println("\x1b[31mError: Failed to generate a password hash\x1b[0m")
        return false
    }

    return bytes.Equal(hash, hashed_password.Hash)
}
