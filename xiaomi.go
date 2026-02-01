package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"time"
)

// miOnOff turns the specified Xiaomi Mi Smart Plug on or off
// Device numbers start at 1. powerOn is true to turn on, false to turn off.
func miOnOff(deviceNum int32, powerOn bool) error {
	devices := MiGetDevices()

	deviceNum = deviceNum + 1

	if deviceNum < 1 || deviceNum > int32(len(devices)) {
		return fmt.Errorf("device number must be between 1 and %d", len(devices))
	}

	device := devices[deviceNum-1]
	token, err := hex.DecodeString(device.Token)
	if err != nil {
		return fmt.Errorf("error decoding token: %v", err)
	}

	// Discover the device
	deviceID, stamp, err := discoverDevice(device.IP)
	if err != nil {
		return fmt.Errorf("discovery error: %v", err)
	}

	// Set power state
	err = setPower(device.IP, token, deviceID, stamp, powerOn)
	if err != nil {
		return fmt.Errorf("failed to set power: %v", err)
	}

	return nil
}

// discoverDevice sends a hello packet to the Xiaomi device and retrieves device ID and timestamp
func discoverDevice(ipAddress string) ([]byte, []byte, error) {
	helloPacket := make([]byte, 32)
	helloPacket[0] = 0x21
	helloPacket[1] = 0x31
	helloPacket[2] = 0x00
	helloPacket[3] = 0x20

	for i := 4; i < 32; i++ {
		helloPacket[i] = 0xFF
	}

	conn, err := net.DialTimeout("udp", fmt.Sprintf("%s:54321", ipAddress), 5*time.Second)
	if err != nil {
		return nil, nil, err
	}
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(5 * time.Second))
	_, err = conn.Write(helloPacket)
	if err != nil {
		return nil, nil, err
	}

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		return nil, nil, err
	}

	if n >= 16 {
		return buffer[8:12], buffer[12:16], nil
	}

	return nil, nil, fmt.Errorf("invalid response")
}

// setPower sends a power command to the Xiaomi device
func setPower(ipAddress string, token []byte, deviceID []byte, stamp []byte, powerOn bool) error {
	state := "off"
	if powerOn {
		state = "on"
	}

	command := map[string]interface{}{
		"id":     1,
		"method": "set_power",
		"params": []interface{}{state},
	}

	jsonData, err := json.Marshal(command)
	if err != nil {
		return err
	}

	encrypted, err := encryptPayload(jsonData, token)
	if err != nil {
		return err
	}

	packet := buildPacket(token, deviceID, stamp, encrypted)

	conn, err := net.DialTimeout("udp", fmt.Sprintf("%s:54321", ipAddress), 5*time.Second)
	if err != nil {
		return err
	}
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(3 * time.Second))
	_, err = conn.Write(packet)
	if err != nil {
		return err
	}

	// Try to read response (may timeout, which is OK)
	buffer := make([]byte, 1024)
	_, _ = conn.Read(buffer)

	return nil
}

// buildPacket constructs a Xiaomi protocol packet with encryption and checksum
func buildPacket(token, deviceID, stamp, encryptedData []byte) []byte {
	packet := make([]byte, 32+len(encryptedData))

	// Magic number and length
	packet[0] = 0x21
	packet[1] = 0x31
	length := uint16(len(packet))
	packet[2] = byte(length >> 8)
	packet[3] = byte(length)

	// Unknown (zeros)
	packet[4] = 0x00
	packet[5] = 0x00
	packet[6] = 0x00
	packet[7] = 0x00

	// Device ID and stamp
	copy(packet[8:12], deviceID)
	copy(packet[12:16], stamp)

	// Encrypted data
	copy(packet[32:], encryptedData)

	// MD5 checksum
	checksumData := make([]byte, 16+len(token)+len(encryptedData))
	copy(checksumData[0:16], packet[0:16])
	copy(checksumData[16:16+len(token)], token)
	copy(checksumData[16+len(token):], encryptedData)
	checksum := md5.Sum(checksumData)
	copy(packet[16:32], checksum[:])

	return packet
}

// encryptPayload encrypts data using AES-CBC with MD5-derived key and IV
func encryptPayload(data []byte, token []byte) ([]byte, error) {
	keyHash := md5.Sum(token)
	key := keyHash[:]

	ivData := append(key, token...)
	ivHash := md5.Sum(ivData)
	iv := ivHash[:]

	padding := aes.BlockSize - len(data)%aes.BlockSize
	padData := make([]byte, len(data)+padding)
	copy(padData, data)
	for i := len(data); i < len(padData); i++ {
		padData[i] = byte(padding)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, len(padData))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, padData)

	return ciphertext, nil
}

// decryptPayload decrypts data using AES-CBC with MD5-derived key and IV
func decryptPayload(encryptedData []byte, token []byte) ([]byte, error) {
	keyHash := md5.Sum(token)
	key := keyHash[:]

	ivData := append(key, token...)
	ivHash := md5.Sum(ivData)
	iv := ivHash[:]

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(encryptedData)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("encrypted data is not a multiple of block size")
	}

	plaintext := make([]byte, len(encryptedData))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plaintext, encryptedData)

	// Remove padding
	if len(plaintext) > 0 {
		padding := int(plaintext[len(plaintext)-1])
		if padding > 0 && padding <= aes.BlockSize && padding <= len(plaintext) {
			plaintext = plaintext[:len(plaintext)-padding]
		}
	}

	return plaintext, nil
}

// miQueryPower queries the actual power state from the Xiaomi device
func miQueryPower(deviceNum int32) (bool, error) {
	devices := MiGetDevices()

	deviceNum = deviceNum + 1

	if deviceNum < 1 || deviceNum > int32(len(devices)) {
		return false, fmt.Errorf("device number must be between 1 and %d", len(devices))
	}

	device := devices[deviceNum-1]
	token, err := hex.DecodeString(device.Token)
	if err != nil {
		return false, fmt.Errorf("error decoding token: %v", err)
	}

	// Discover the device
	deviceID, stamp, err := discoverDevice(device.IP)
	if err != nil {
		return false, fmt.Errorf("discovery error: %v", err)
	}

	// Query power state
	command := map[string]interface{}{
		"id":     1,
		"method": "get_prop",
		"params": []interface{}{"power"},
	}

	jsonData, err := json.Marshal(command)
	if err != nil {
		return false, err
	}

	encrypted, err := encryptPayload(jsonData, token)
	if err != nil {
		return false, err
	}

	packet := buildPacket(token, deviceID, stamp, encrypted)

	conn, err := net.DialTimeout("udp", fmt.Sprintf("%s:54321", device.IP), 5*time.Second)
	if err != nil {
		return false, err
	}
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(3 * time.Second))
	_, err = conn.Write(packet)
	if err != nil {
		return false, err
	}

	// Read response
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		return false, fmt.Errorf("failed to read response: %v", err)
	}

	if n < 32 {
		return false, fmt.Errorf("response too short")
	}

	// Decrypt the response payload
	encryptedResponse := buffer[32:n]
	decrypted, err := decryptPayload(encryptedResponse, token)
	if err != nil {
		return false, fmt.Errorf("failed to decrypt response: %v", err)
	}

	// Parse JSON response
	var response struct {
		Result []string `json:"result"`
	}
	err = json.Unmarshal(decrypted, &response)
	if err != nil {
		return false, fmt.Errorf("failed to parse response: %v", err)
	}

	if len(response.Result) > 0 {
		return response.Result[0] == "on", nil
	}

	return false, fmt.Errorf("no power state in response")
}
