package models

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

// Protocol16Type representa los tipos de datos en el protocolo.
var Protocol16Type = map[string]byte{
	"Unknown":           0,
	"Null":              42,
	"Byte":              68, //68
	"Boolean":           97,
	"Short":             98,
	"Integer":           100,
	"IntegerArray":      101,
	"Double":            102,
	"Long":              105,
	"Float":             104,
	"String":            107,
	"StringArray":       108,
	"ByteArray":         110,
	"EventData":         111,
	"Dictionary":        112,
	"Array":             113,
	"OperationResponse": 115,
	"OperationRequest":  120,
	"Hashtable":         121,
	"ObjectArray":       122,
}

// Deserialize deserializa los datos de entrada según el tipo de código.
func Deserialize(input *bytes.Buffer, typeCode byte) (interface{}, error) {

	switch typeCode {
	case Protocol16Type["Unknown"], Protocol16Type["Null"]:
		return nil, nil
	case Protocol16Type["Byte"]:
		return DeserializeByte(input)
	case Protocol16Type["Boolean"]:
		return DeserializeBoolean(input)
	case Protocol16Type["Short"]:
		return DeserializeShort(input)
	case Protocol16Type["Integer"]:
		return DeserializeInteger(input)
	case Protocol16Type["IntegerArray"]:
		return DeserializeIntegerArray(input)
	case Protocol16Type["Double"]:
		return DeserializeDouble(input)
	case Protocol16Type["Long"]:
		return DeserializeLong(input)
	case Protocol16Type["Float"]:
		return DeserializeFloat(input)
	case Protocol16Type["String"]:
		return DeserializeString(input)
	case Protocol16Type["StringArray"]:
		return DeserializeStringArray(input)
	case Protocol16Type["ByteArray"]:
		return DeserializeByteArray(input)
	case Protocol16Type["EventData"]:
		return DeserializeEventData(input)
	case Protocol16Type["Dictionary"]:
		return DeserializeDictionary(input)
	case Protocol16Type["Array"]:
		return DeserializeArray(input)
	case Protocol16Type["OperationResponse"]:
		return DeserializeOperationResponse(input)
	case Protocol16Type["OperationRequest"]:
		return DeserializeOperationRequest(input)
	case Protocol16Type["Hashtable"]:
		return DeserializeHashtable(input)
	case Protocol16Type["ObjectArray"]:
		return DeserializeObjectArray(input)
	default:
		// fmt.Println("el codigo: ' ", typeCode, " ' no se encuentra")
		return nil, errors.New("type code not implemented")
	}
}

// DeserializeByte deserializa un byte.
func DeserializeByte(input *bytes.Buffer) (byte, error) {
	return input.ReadByte()
}

// DeserializeBoolean deserializa un booleano.
func DeserializeBoolean(input *bytes.Buffer) (bool, error) {
	val, err := input.ReadByte()
	return val != 0, err
}

// DeserializeShort deserializa un short.
func DeserializeShort(input *bytes.Buffer) (uint16, error) {
	var shortVal uint16
	err := binary.Read(input, binary.BigEndian, &shortVal)
	return shortVal, err
}

// DeserializeInteger deserializa un entero.
func DeserializeInteger(input *bytes.Buffer) (uint32, error) {
	var intVal uint32
	err := binary.Read(input, binary.BigEndian, &intVal)
	return intVal, err
}

func DeserializeIntegerArray(input *bytes.Buffer) ([]uint32, error) {
	size, err := DeserializeInteger(input)
	var res []uint32
	var val uint32
	for i := 0; i < int(size); i++ {
		val, err = DeserializeInteger(input)
		res = append(res, val)
	}

	return res, err
}

func DeserializeDouble(input *bytes.Buffer) (float64, error) {
	var doubleVal float64
	err := binary.Read(input, binary.BigEndian, &doubleVal)
	return doubleVal, err
}

func DeserializeLong(input *bytes.Buffer) (int64, error) {
	var longVal int64
	err := binary.Read(input, binary.BigEndian, &longVal)
	if err != nil {
		return 0, err
	}

	return longVal, nil
}

func DeserializeFloat(input *bytes.Buffer) (float32, error) {
	var floatVal float32
	err := binary.Read(input, binary.BigEndian, &floatVal)
	return floatVal, err
}

// DeserializeString deserializa una cadena de caracteres desde un buffer de bytes.
func DeserializeString(input *bytes.Buffer) (string, error) {
	stringSize, err := DeserializeShort(input)
	if err != nil {
		return "", err
	}
	if stringSize == 0 {
		return "", nil
	}

	res := make([]byte, stringSize)
	if _, err := input.Read(res); err != nil {
		return "", err
	}

	return string(res), nil
}

// DeserializeByteArray deserializa un slice de bytes desde un buffer de bytes.
func DeserializeByteArray(input *bytes.Buffer) ([]byte, error) {
	arraySize, err := DeserializeInteger(input)
	if err != nil {
		return nil, err
	}

	res := make([]byte, arraySize)
	if _, err := input.Read(res); err != nil {
		return nil, err
	}

	return res, nil
}

// DeserializeArray deserializa un slice de interface{} desde un buffer de bytes.
func DeserializeArray(input *bytes.Buffer) ([]interface{}, error) {
	size, err := DeserializeShort(input)
	if err != nil {
		return nil, err
	}

	typeCode, err := DeserializeByte(input)
	if err != nil {
		return nil, err
	}

	res := make([]interface{}, size)
	for i := uint16(0); i < size; i++ {
		val, err := Deserialize(input, typeCode)
		if err != nil {
			return nil, err
		}
		res[i] = val
	}

	return res, nil
}

// DeserializeStringArray deserializa un slice de cadenas de caracteres desde un buffer de bytes.
func DeserializeStringArray(input *bytes.Buffer) ([]string, error) {
	size, err := DeserializeShort(input)
	if err != nil {
		return nil, err
	}

	res := make([]string, size)
	for i := uint16(0); i < size; i++ {
		val, err := DeserializeString(input)
		if err != nil {
			return nil, err
		}
		res[i] = val
	}

	return res, nil
}

// DeserializeObjectArray deserializa un slice de interface{} desde un buffer de bytes.
func DeserializeObjectArray(input *bytes.Buffer) ([]interface{}, error) {
	tableSize, err := DeserializeShort(input)
	if err != nil {
		return nil, err
	}

	output := make([]interface{}, tableSize)
	for i := uint16(0); i < tableSize; i++ {
		typeCode, err := DeserializeByte(input)
		if err != nil {
			return nil, err
		}
		val, err := Deserialize(input, typeCode)
		if err != nil {
			return nil, err
		}
		output[i] = val
	}

	return output, nil
}

// DeserializeHashtable deserializa una tabla hash desde un buffer de bytes.
func DeserializeHashtable(input *bytes.Buffer) (map[interface{}]interface{}, error) {
	tableSize, err := DeserializeShort(input)
	if err != nil {
		return nil, err
	}

	return DeserializeDictionaryElements(input, int(tableSize), 0, 0)
}

// DeserializeDictionary deserializa un diccionario desde un buffer de bytes.
func DeserializeDictionary(input *bytes.Buffer) (map[interface{}]interface{}, error) {
	keyTypeCode, err := DeserializeByte(input)
	if err != nil {
		return nil, err
	}
	valueTypeCode, err := DeserializeByte(input)
	if err != nil {
		return nil, err
	}
	dictionarySize, err := DeserializeShort(input)
	if err != nil {
		return nil, err
	}

	return DeserializeDictionaryElements(input, int(dictionarySize), keyTypeCode, valueTypeCode)
}

// DeserializeDictionaryElements deserializa los elementos de un diccionario.
func DeserializeDictionaryElements(input *bytes.Buffer, dictionarySize int, keyTypeCode byte, valueTypeCode byte) (map[interface{}]interface{}, error) {
	output := make(map[interface{}]interface{})

	for i := 0; i < dictionarySize; i++ {
		var key interface{}
		var value interface{}
		var err error

		// Determinar el tipo de clave si es necesario.
		if keyTypeCode == 0 || keyTypeCode == 42 {
			keyTypeCode, err = DeserializeByte(input)
			if err != nil {
				return nil, err
			}
		}
		key, err = Deserialize(input, keyTypeCode)
		if err != nil {
			return nil, err
		}

		// Determinar el tipo de valor si es necesario.
		if valueTypeCode == 0 || valueTypeCode == 42 {
			valueTypeCode, err = DeserializeByte(input)
			if err != nil {
				return nil, err
			}
		}
		value, err = Deserialize(input, valueTypeCode)
		if err != nil {
			return nil, err
		}

		output[key] = value
	}

	return output, nil
}

// DeserializeOperationRequest deserializa una solicitud de operación desde un buffer de bytes.
func DeserializeOperationRequest(input *bytes.Buffer) (map[string]interface{}, error) {
	operationCode, err := DeserializeByte(input)
	if err != nil {
		return nil, err
	}
	parameters, err := DeserializeParameterTable(input)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"operationCode": operationCode,
		"parameters":    parameters,
	}, nil
}

// DeserializeOperationResponse deserializa una respuesta de operación desde un buffer de bytes.
func DeserializeOperationResponse(input *bytes.Buffer) (map[string]interface{}, error) {
	operationCode, err := DeserializeByte(input)
	if err != nil {
		return nil, err
	}
	returnCode, err := DeserializeShort(input)
	if err != nil {
		return nil, err
	}
	debugMessageTypeCode, err := DeserializeByte(input)
	if err != nil {
		return nil, err
	}
	debugMessage, err := Deserialize(input, debugMessageTypeCode)
	if err != nil {
		return nil, err
	}
	parameters, err := DeserializeParameterTable(input)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"operationCode": operationCode,
		"returnCode":    returnCode,
		"debugMessage":  debugMessage,
		"parameters":    parameters,
	}, nil
}

// DeserializeEventData deserializa datos de eventos desde un buffer de bytes.
func DeserializeEventData(input *bytes.Buffer) (map[string]interface{}, error) {
	code, err := DeserializeByte(input)
	if err != nil {

		return nil, err
	}

	parameters, err := DeserializeParameterTable(input)
	if err != nil {
		return nil, err
	}

	if code == 3 {
		byteArray, ok := parameters[1].([]byte)
		if !ok {
			return nil, fmt.Errorf("no se encontro parametro")
		}
		reader := bytes.NewReader(byteArray)
		reader.Seek(9, io.SeekStart)
		var position0 float64

		err = binary.Read(reader, binary.LittleEndian, &position0)

		if err != nil {
			return nil, err
		}

		var position1 float64

		err = binary.Read(reader, binary.LittleEndian, &position1)
		if err != nil {
			return nil, err
		}
		parameters[4] = position0
		parameters[5] = position1
		parameters[252] = byte(3)
	}

	// Aquí deberías agregar la lógica específica para manejar el código 3 y modificar los parámetros como en tu función de JavaScript.

	return map[string]interface{}{
		"code":       code,
		"parameters": parameters,
	}, nil
}

// DeserializeParameterTable deserializa una tabla de parámetros desde un buffer de bytes.
func DeserializeParameterTable(input *bytes.Buffer) (map[byte]interface{}, error) {
	// Asumiendo que la función DeserializeShort ya ha sido implementada.
	tableSize, err := DeserializeShort(input)
	if err != nil {
		return nil, err
	}

	table := make(map[byte]interface{})

	for i := 0; i < int(tableSize); i++ {
		key, err := DeserializeByte(input)
		if err != nil {
			return nil, err
		}

		valueTypeCode, err := DeserializeByte(input)
		if err != nil {
			return nil, err
		}

		value, err := Deserialize(input, valueTypeCode)
		if err != nil {

			return nil, err
		}

		table[key] = value
	}

	return table, nil
}
