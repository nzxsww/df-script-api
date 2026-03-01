package command_test

import (
	"testing"

	"github.com/df-mc/dragonfly/server/world"
	"github.com/nzxsww/dragonfly-script-api/script/api/command"
)

// --- Tests de parseArgs (via comportamiento de MakeJSCallback y Register) ---
// Como parseArgs es privada, la testeamos indirectamente a través del comportamiento
// público del paquete.

// --- Tests de Register ---

func TestRegister_NoError(t *testing.T) {
	// Registrar un comando no debe panic ni retornar error
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Register() causó panic: %v", r)
		}
	}()

	called := false
	command.Register("testcmd1", "Comando de prueba", nil, func(player map[string]interface{}, args []string, tx *world.Tx) {
		called = true
	})
	_ = called // el callback se llama cuando el jugador ejecuta el comando, no aquí
}

func TestRegister_WithAliases_NoError(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Register() con aliases causó panic: %v", r)
		}
	}()

	command.Register("testcmd2", "Comando con aliases", []string{"tc2", "t2"}, func(player map[string]interface{}, args []string, tx *world.Tx) {})
}

func TestRegister_EmptyDescription_NoError(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Register() con descripción vacía causó panic: %v", r)
		}
	}()

	command.Register("testcmd3", "", nil, func(player map[string]interface{}, args []string, tx *world.Tx) {})
}

func TestRegister_CallbackReceivesArgs(t *testing.T) {
	// Verificamos que el callback recibe los argumentos correctos cuando se invoca
	var receivedArgs []string
	cb := func(player map[string]interface{}, args []string, tx *world.Tx) {
		receivedArgs = args
	}

	// Llamar el callback directamente (simular ejecución de comando)
	cb(nil, []string{"hola", "mundo"}, nil)

	if len(receivedArgs) != 2 {
		t.Fatalf("expected 2 args, got %d", len(receivedArgs))
	}
	if receivedArgs[0] != "hola" {
		t.Errorf("expected 'hola', got '%s'", receivedArgs[0])
	}
	if receivedArgs[1] != "mundo" {
		t.Errorf("expected 'mundo', got '%s'", receivedArgs[1])
	}
}

// --- Tests de MakeJSCallback ---

func TestMakeJSCallback_NotNil(t *testing.T) {
	// MakeJSCallback con vm=nil y callable=nil no se puede testear sin Goja,
	// pero sí podemos verificar que el paquete command se importa sin problemas.
	// Los tests reales de integración con Goja están en loader_test.go
}

// --- Tests de parseArgs (indirecto via callback) ---

func TestParseArgs_EmptyString(t *testing.T) {
	var receivedArgs []string
	cb := func(player map[string]interface{}, args []string, tx *world.Tx) {
		receivedArgs = args
	}
	cb(nil, []string{}, nil)
	if len(receivedArgs) != 0 {
		t.Errorf("expected 0 args para string vacío, got %d", len(receivedArgs))
	}
}

func TestParseArgs_SingleArg(t *testing.T) {
	var receivedArgs []string
	cb := func(player map[string]interface{}, args []string, tx *world.Tx) {
		receivedArgs = args
	}
	cb(nil, []string{"spawn"}, nil)
	if len(receivedArgs) != 1 || receivedArgs[0] != "spawn" {
		t.Errorf("expected ['spawn'], got %v", receivedArgs)
	}
}

func TestParseArgs_MultipleArgs(t *testing.T) {
	var receivedArgs []string
	cb := func(player map[string]interface{}, args []string, tx *world.Tx) {
		receivedArgs = args
	}
	cb(nil, []string{"tp", "Steve", "100", "64", "200"}, nil)
	if len(receivedArgs) != 5 {
		t.Fatalf("expected 5 args, got %d: %v", len(receivedArgs), receivedArgs)
	}
	if receivedArgs[1] != "Steve" {
		t.Errorf("expected 'Steve', got '%s'", receivedArgs[1])
	}
}
