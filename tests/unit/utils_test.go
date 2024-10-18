package unit

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/suidevv/tableye-api/utils"
)

const (
	privateKey = `LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlDV3dJQkFBS0JnUURLc1UzNUNiVWQ1cXZHbGpEcTVZWHIxZlNybHZTTHprMDFUWDhGamdhWnRCdW1VdW1NClQyZjZqNU1EbXp0bEhqSXBrWm4rdk42RzZ6bGhWYXNONEdFZVFhUXNXS0NWSzRUODF0L1dtRFkwQUJramk1akQKY25WVEFUOWRoV25neDUwUC95TFFPYWJIMWd2akRxQzJEeXVFaUFnWDZocURjKzhmRm9WVlp3Vm5ld0lEQVFBQgpBb0dBS1diWENPL3BMVjdoSm5LbU1reklzcGZrM3FtNmNOWW1ZaVZldFRsQjh1SmRwWGNaR2w1YjNFdTRXVXU3CmNaZWQ0bXpKdWtWRTVPVW1OdElEV3hYQ2NFaTB0ZExob3hXSWJFOGtvSHo5QTBiekRpcXJnd3V0LzFIU1MrbDMKUGNhUU5BZXp0dENDY0V3UWdkd21NcFpKTWE5QlBjK3orV3FEUUFoK2dtMHhXd2tDUVFENjY5RTA2ZmxyRWFMbQpWSTAxOGNuYXJlY2lzdEVocmxFaWM4WVFxandzVWZveTFCMnVIU3gyRGdTQ0haTjk1WEdZZ3dmOGdJNkg5UFR2CnpQNVV1TXZEQWtFQXpzdVZtVnFSQlNTWDNHMVlqbXlIK3R5NVBzakp6Rmc1TlZQMzJDRzd6clU0VzBNbGVYaG0KTmYxNEhBdGVPV3htQlF1Q3EvR3h0aTNPdE5aZHJNUVI2UUpBYWZwdndnbVFic2hrSlNSUkFCZS9TYjFwZ2g1RQpkaFZKNzJNMnBKTkNGdllJMXE4QVdpbTRQYVJ1QXdhNjVOR2p5T2FPMlBielBEa1p1cTY2UE01UVFRSkFEYlF5CkVycU10N0dJR3NSb1JPL3VSdktQbUJpSVB2RnR3Um55WjdFOGwrTXNlK2ZFT1B1QWtuMWNrMGN4bEU2WnFDWHUKSCtUaGFQZzZKWU83SzNMRzJRSkFCZCtZQVFJRWVXeXNPSWF0TXNSUWlqdnd0MHJzNjVHNGQvVWd4K2tiVENQaQpkYmNVbzB3Y1hwQU9lb1ZnSWQrR2hpMUo3YzczV0NtL3RCY1lWWWhVWEE9PQotLS0tLUVORCBSU0EgUFJJVkFURSBLRVktLS0tLQo=`
	publicKey  = `LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUlHZk1BMEdDU3FHU0liM0RRRUJBUVVBQTRHTkFEQ0JpUUtCZ1FES3NVMzVDYlVkNXF2R2xqRHE1WVhyMWZTcgpsdlNMemswMVRYOEZqZ2FadEJ1bVV1bU1UMmY2ajVNRG16dGxIaklwa1puK3ZONkc2emxoVmFzTjRHRWVRYVFzCldLQ1ZLNFQ4MXQvV21EWTBBQmtqaTVqRGNuVlRBVDlkaFduZ3g1MFAveUxRT2FiSDFndmpEcUMyRHl1RWlBZ1gKNmhxRGMrOGZGb1ZWWndWbmV3SURBUUFCCi0tLS0tRU5EIFBVQkxJQyBLRVktLS0tLQo=`
)

func TestCreateToken(t *testing.T) {
	payload := "testuser"
	duration := time.Minute * 15
	token, err := utils.CreateToken(duration, payload, privateKey)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Add assertions to verify token contents
	claims, err := utils.ValidateToken(token, publicKey)
	assert.NoError(t, err)
	assert.Equal(t, payload, claims)
	// Add more assertions as needed
}

func TestValidateToken(t *testing.T) {
	payload := "testuser"
	duration := time.Minute * 15
	token, err := utils.CreateToken(duration, payload, privateKey)
	assert.NoError(t, err)

	// Test with valid token
	claims, err := utils.ValidateToken(token, publicKey)
	assert.NoError(t, err)
	assert.Equal(t, payload, claims)

	// Test with invalid token
	_, err = utils.ValidateToken("invalid*token", publicKey)
	assert.Error(t, err)

	// Test with expired token
	expiredToken, err := utils.CreateToken(-time.Minute, payload, privateKey)
	assert.NoError(t, err)
	_, err = utils.ValidateToken(expiredToken, publicKey)
	assert.Error(t, err)

	// Test with token signed by different private key
	differentPrivateKey := `LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlDWEFJQkFBS0JnUUN1R01zUkZKRlZkMmZGUWlubzRKeGhtU3hOQzdneHUwY1JLREMwN3JUbVVaMUpBTkxMCjZQWDlWQkxlVW1kUE93WWtzaGZNTERuTndRUW1jVkdTRDFTVm1OV09KODlBQzlWdVZzVjVXMXhkNnliNVJkbXgKOFdkVks4TGlrTzRUVHBNVFAyQ3Rpd1pwMDA4VFBubjM5ak9RMHE2aUpycXVxUUJybE42ejg3V25DUUlEQVFBQgpBb0dBS2Fxclp5SWFDbTl5ZVlaMVMxUCtlV2xVYmZnaEZGM3pDMHNvSlJXQVhDTXFWcCtJMDk2a2lRWEhJa2hzCmkwemdFb1pCSmM5cjRpK1EwN2FKUkVEQkNFdEVLUXU1MkR6dXFtQnhSb2R4alhlaDB3dGN1QUtjeWQzU3hlbU8KLzJ0NkRsN3hTZlMrRXVpNDgvdDNUQnZkcDh1UHFVSkF3NUNmQThmaXNZREp4eFVDUVFEbWFMdVlXMWpIRXMraApZajRLTEg0NFRpbHRpVUovVGRZL2NDV2t5NGUzUm5wWEoySDRWUko3enVHc0Z0S2NDVEtud1NtVXR1VEFoa0VBCk1wQUZ1KzFEQWtFQXdXN3NRSUllcnVJUTdHVnJVcDg0S0QzUVJDdGRELzdScFYwcjcwVWxsMTkvOTlsUHB5QmwKck4wUkdoUUc3VFZLMXZhWStEZjBZdUN3dWpRbjY2TVB3d0pBQitBTWlXaVY0RGdFWUwrNjN4NG1Na1o1cEFUTgpBUXpvQmNNUGhsSnVrUlVYbVdML05qMnlKQWt1TFhPYVB6c1JRQ3FhQVRzL0ZsV0FZMEZYS3RzQmdRSkFidkhGCkZaYk1MSGhEUnFOQTlDbVlWeFJsSU1SU1l6cy9XWDVnRmFOdVZTMFVRNzdqZmJNS1BpU3BpM0NUTEhpVmpVZngKSXVWTkNXMWdUOXhjVFQzQWF3SkJBSkFFWXFHNkcrdTRMRGc2Rzk4cmRzMHc2UE9ESDNLRXBoTkRWeHh6M2hTcQo2SW85WnExeXlJZ3pqQmZIRTdRbm9UMW5GM0d1ZzlyakE1eGVuNjRBKzNFPQotLS0tLUVORCBSU0EgUFJJVkFURSBLRVktLS0tLQo=`
	differentToken, err := utils.CreateToken(duration, payload, differentPrivateKey)
	assert.NoError(t, err)
	_, err = utils.ValidateToken(differentToken, publicKey)
	assert.Error(t, err)
}

func TestHashPassword(t *testing.T) {
	password := "testpassword123"

	hashedPassword, err := utils.HashPassword(password)
	assert.NoError(t, err)
	assert.NotEqual(t, password, hashedPassword)

	err = utils.VerifyPassword(hashedPassword, password)
	assert.NoError(t, err)

	err = utils.VerifyPassword(hashedPassword, "wrongpassword")
	assert.Error(t, err)
}

func TestVerifyPassword(t *testing.T) {
	password := "testpassword123"
	hashedPassword, _ := utils.HashPassword(password)

	err := utils.VerifyPassword(hashedPassword, password)
	assert.NoError(t, err)

	err = utils.VerifyPassword(hashedPassword, "wrongpassword")
	assert.Error(t, err)
}
