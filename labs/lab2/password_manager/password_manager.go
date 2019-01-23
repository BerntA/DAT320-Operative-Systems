// +build !solution

package main

import (
	"bufio"
	"os"
	"strings"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"errors"
	"io/ioutil"
	"encoding/json"
)

/*
Task 7: Password Manager

This task focuses on building a password manager that stores
fictional passwords for websites. These passwords should be
handled somewhat securely, which is why only a hash will be stored.

The password manager should be implemented through a command line
application; allowing the user to execute all the functions in the
PasswordManager interface.

*/

// PasswordManager : Holds passwords for certain websites
type PasswordManager struct {
	data map[string][]byte
}

// NewPasswordManager returns an initialized instance of PasswordManager
func NewPasswordManager() *PasswordManager {
	var pwdmngr PasswordManager
	pwdmngr.data = make (map[string][]byte)
	return &pwdmngr
}

// Set creates or updates the password associated with the given site.
// The stored password should be hashed and salted, which can be accomplished
// with the bcrypt package (https://godoc.org/golang.org/x/crypto/bcrypt).
// Use the command "go get golang.org/x/crypto/bcrypt" if you do not have it installed.
func (m *PasswordManager) Set(site, password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err == nil {
		m.data[site] = bytes
	}
	return err
}

// Get returns the stored hash for the given site.
func (m *PasswordManager) Get(site string) []byte {
	value, ok := m.data[site]
	if ok {
		return value
	}

	return nil
}

// Verify checks whether the password given for a site matches the stored password.
// If the returned error value is nil, the passwords match.
// Hint: The bcrypt package may be of use here
func (m *PasswordManager) Verify(site, password string) error {
	value, ok := m.data[site]
	if ok {
		return bcrypt.CompareHashAndPassword(value, []byte(password))
	}

	msg := fmt.Sprintf("Unable to find website %v", site)
	return errors.New(msg)
}

// Remove deletes the password for a given site from the password manager.
func (m *PasswordManager) Remove(site string) {
	_, ok := m.data[site]
	if ok {
		delete(m.data, site)
	}
}

// Save stores all the passwords in a file of the given name.
// Read the ioutil documentation for inforamtion on reading/writing files
// https://golang.org/pkg/io/ioutil/
//
// The file contents should be serialized in some way, for example using JSON
// (https://golang.org/pkg/encoding/json/) or XML (https://golang.org/pkg/encoding/xml/) etc.
func (m *PasswordManager) Save(fileName string) error {
	if len(m.data) <= 0 {
		return errors.New("No passwords found")
	}

	b, err := json.MarshalIndent(&m.data, "", "\t")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(fileName, b, 0644)
}

// Load reads the given file and decodes the values to replace the state of the
// PasswordManager.
func (m *PasswordManager) Load(fileName string) error {
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}

	output := make(map[string][]byte)
	err = json.Unmarshal(bytes, &output)
	
	if err == nil {
		m.data = output
	}

	return err
}

// Encrypt encrypts the input data using a cipher of your choice.
// This should be performed before saving passwords to disk.
// Go supports the AES(https://golang.org/pkg/crypto/aes/) and
// DES(https://golang.org/pkg/crypto/des/) ciphers. Note that
// TripleDES should be used if you choose DES.
//
// https://golang.org/pkg/crypto/cipher contains examples of different
// cipher modes. The simplest mode is block mode, although you are free
// to use stream mode if you prefer.
//
// It is also important to use keys of the right size depending on which
// cipher you chose. For example the AES key size is 16 bytes, or regular
// latin characters. You can either force the given key to 16 bytes by repeating
// or cutting off the key, or you can just rely on the user inputting 16 bytes.
func (m *PasswordManager) Encrypt(plaintext []byte, keystring string) []byte {
	// NOTE: This implementation is OPTIONAL
	return nil
}

// Decrypt reverses the operations performed by Encrypt, and is
// used when loading the password file from storage.
func (m *PasswordManager) Decrypt(ciphertext []byte, keystring string) []byte {
	// NOTE: This implementation is OPTIONAL
	return nil
}

// This is the main function of the application.
// User input should be continuously read and checked for commands
// for all the defined operations.
// See https://golang.org/pkg/bufio/#Reader and especially the ReadLine
// function.
func main() {
	fmt.Println(
`Password Manager has started!

Commands:
create
get
verify
remove
save
load

quit to quit
`)

	manager := NewPasswordManager()

	scan := bufio.NewScanner(os.Stdin) 
	for scan.Scan() {
		cmd := strings.ToLower(scan.Text())
		if strings.HasPrefix(cmd, "quit") {
			fmt.Println("Quitting!")
			break
		}

		if strings.HasPrefix(cmd, "create") {
			data := strings.Split(cmd, " ")
			if len(data) != 3 {
				fmt.Println("create requires 2 arguments: site + password")
				continue
			}

			err := manager.Set(data[1], data[2])
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Printf("Successfully created password %v for site %v!\n", data[2], data[1])
			}
		} else if strings.HasPrefix(cmd, "get") {
			data := strings.Split(cmd, " ")
			if len(data) != 2 {
				fmt.Println("get requires 1 argument: site")
				continue
			}

			v := manager.Get(data[1])
			if v != nil {
				fmt.Printf("Found password %v for site %v!\n", string(v[:]), data[1])
			} else {
				fmt.Printf("Site %v wasn't found!\n", data[1])
			}
		} else if strings.HasPrefix(cmd, "verify") {
			data := strings.Split(cmd, " ")
			if len(data) != 3 {
				fmt.Println("verify requires 2 arguments: site + password")
				continue
			}

			err := manager.Verify(data[1], data[2])
			if err != nil {
				fmt.Printf("Failed to verify password %v for site %v!\n", data[2], data[1])
			} else {
				fmt.Printf("Successfully verified password %v for site %v!\n", data[2], data[1])
			}
		} else if strings.HasPrefix(cmd, "remove") {
			data := strings.Split(cmd, " ")
			if len(data) != 2 {
				fmt.Println("remove requires 1 argument: site")
				continue
			}

			manager.Remove(data[1])
			fmt.Printf("Tried to remove saved password for site %v\n", data[1])
		} else if strings.HasPrefix(cmd, "save") {
			data := strings.Split(cmd, " ")
			if len(data) != 2 {
				fmt.Println("save requires 1 argument: filename")
				continue
			}

			err := manager.Save(data[1])
			if err != nil {
				fmt.Printf("(%v) Unable to save passwords to '%v'!\n", err.Error(), data[1])
			} else {
				fmt.Printf("Saved passwords successfully to '%v'\n", data[1])
			}
		} else if strings.HasPrefix(cmd, "load") {
			data := strings.Split(cmd, " ")
			if len(data) != 2 {
				fmt.Println("load requires 1 argument: filename")
				continue
			}	

			err := manager.Load(data[1])
			if err != nil {
				fmt.Printf("(%v) Unable to load passwords from '%v'!\n", err.Error(), data[1])
			} else {
				fmt.Printf("Successfully loaded passwords from '%v'\n", data[1])
			}
		}
	}
}
