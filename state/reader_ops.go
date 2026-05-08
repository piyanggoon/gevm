package state

import (
	"fmt"
	"reflect"

	"github.com/holiman/uint256"

	"github.com/Giulio2002/gevm/types"
)

// ReaderOps describes how to call an external state reader whose concrete
// address/account types live outside GEVM.
type ReaderOps struct {
	Address    func(types.Address) any
	StorageKey func(uint256.Int) any
	Account    func(any, types.Bytes) (AccountInfo, bool, error)
	Basic      func(types.Address) (AccountInfo, bool, error)
	Storage    func(types.Address, uint256.Int) (uint256.Int, error)
	HasStorage func(types.Address) (bool, error)
	Code       func(types.Address) (types.Bytes, error)
	BlockHash  func(uint64) (types.B256, error)
}

func (j *Journal) Basic(address types.Address) (AccountInfo, bool, error) {
	if db, ok := j.DB.(Database); ok {
		return db.Basic(address)
	}
	if j.ReaderOps.Basic != nil {
		return j.ReaderOps.Basic(address)
	}
	if j.DB == nil {
		return AccountInfo{}, false, nil
	}
	if values, ok, err := j.callReaderWith("ReadAccountDataRaw", func(t reflect.Type) (reflect.Value, error) {
		return addressArg(address, t)
	}); ok {
		if err != nil {
			return AccountInfo{}, false, err
		}
		if err := resultError(values[len(values)-1]); err != nil {
			return AccountInfo{}, false, err
		}
		acc := resultInterface(values[len(values)-2])
		if acc == nil {
			return AccountInfo{}, false, nil
		}
		info, exists, err := accountInfo(acc, nil)
		if err != nil || !exists {
			return info, exists, err
		}
		code, err := j.readCode(address)
		if err != nil {
			return AccountInfo{}, false, err
		}
		if len(code) > 0 {
			info.Code = types.BytesFrom(code)
		}
		return info, true, nil
	}
	addr, err := j.readerAddress(address)
	if err != nil {
		return AccountInfo{}, false, err
	}
	acc, err := j.callReader1("ReadAccountData", addr)
	if err != nil {
		return AccountInfo{}, false, err
	}
	if acc == nil {
		return AccountInfo{}, false, nil
	}
	if j.ReaderOps.Account == nil {
		return accountInfo(acc, nil)
	}
	info, exists, err := j.ReaderOps.Account(acc, nil)
	if err != nil || !exists {
		return info, exists, err
	}
	code, err := j.readCode(addr)
	if err != nil {
		return AccountInfo{}, false, err
	}
	if len(code) > 0 {
		info.Code = types.BytesFrom(code)
	}
	return info, true, nil
}

func (j *Journal) CodeByHash(codeHash types.B256) (types.Bytes, error) {
	if db, ok := j.DB.(Database); ok {
		return db.CodeByHash(codeHash)
	}
	return nil, fmt.Errorf("code hash %s was not loaded with its account", codeHash.Hex())
}

func (j *Journal) ReadCode(address types.Address) (types.Bytes, error) {
	return j.readCode(address)
}

func (j *Journal) Storage(address types.Address, index uint256.Int) (uint256.Int, error) {
	if db, ok := j.DB.(Database); ok {
		return db.Storage(address, index)
	}
	if j.ReaderOps.Storage != nil {
		return j.ReaderOps.Storage(address, index)
	}
	if j.DB == nil {
		return uint256.Int{}, nil
	}
	if values, ok, err := j.callReaderWith("ReadAccountStorageRaw", func(t reflect.Type) (reflect.Value, error) {
		if t.Size() == 20 {
			return addressArg(address, t)
		}
		return storageKeyArg(index, t)
	}); ok {
		if err != nil {
			return uint256.Int{}, err
		}
		if err := resultError(values[len(values)-1]); err != nil {
			return uint256.Int{}, err
		}
		if value := resultInterface(values[0]); value != nil {
			return value.(uint256.Int), nil
		}
		return uint256.Int{}, nil
	}
	addr, err := j.readerAddress(address)
	if err != nil {
		return uint256.Int{}, err
	}
	key, err := j.readerStorageKey(index)
	if err != nil {
		return uint256.Int{}, err
	}
	value, err := j.callReader1("ReadAccountStorage", addr, key)
	if err != nil {
		return uint256.Int{}, err
	}
	if value == nil {
		return uint256.Int{}, nil
	}
	return value.(uint256.Int), nil
}

func (j *Journal) HasStorage(address types.Address) (bool, error) {
	if db, ok := j.DB.(Database); ok {
		return db.HasStorage(address)
	}
	if j.ReaderOps.HasStorage != nil {
		return j.ReaderOps.HasStorage(address)
	}
	if j.DB == nil {
		return false, nil
	}
	if values, ok, err := j.callReaderWith("HasStorageRaw", func(t reflect.Type) (reflect.Value, error) {
		return addressArg(address, t)
	}); ok {
		if err != nil {
			return false, err
		}
		if err := resultError(values[len(values)-1]); err != nil {
			return false, err
		}
		return resultInterface(values[0]).(bool), nil
	}
	addr, err := j.readerAddress(address)
	if err != nil {
		return false, err
	}
	value, err := j.callReader1("HasStorage", addr)
	if err != nil {
		return false, err
	}
	return value.(bool), nil
}

func (j *Journal) BlockHash(number uint64) (types.B256, error) {
	if db, ok := j.DB.(Database); ok {
		return db.BlockHash(number)
	}
	if j.ReaderOps.BlockHash == nil {
		return types.B256Zero, nil
	}
	return j.ReaderOps.BlockHash(number)
}

func (j *Journal) readerAddress(address types.Address) (any, error) {
	if j.ReaderOps.Address == nil {
		return nil, fmt.Errorf("external reader address converter is not configured")
	}
	return j.ReaderOps.Address(address), nil
}

func (j *Journal) readerStorageKey(index uint256.Int) (any, error) {
	if j.ReaderOps.StorageKey == nil {
		return nil, fmt.Errorf("external reader storage-key converter is not configured")
	}
	return j.ReaderOps.StorageKey(index), nil
}

func (j *Journal) readCode(address any) (types.Bytes, error) {
	if j.ReaderOps.Code != nil {
		addr, ok := address.(types.Address)
		if ok {
			return j.ReaderOps.Code(addr)
		}
	}
	if addr, ok := address.(types.Address); ok && j.DB != nil {
		if values, found, err := j.callReaderWith("ReadAccountCodeRaw", func(t reflect.Type) (reflect.Value, error) {
			return addressArg(addr, t)
		}); found {
			if err != nil {
				return nil, err
			}
			if err := resultError(values[len(values)-1]); err != nil {
				return nil, err
			}
			if code := resultInterface(values[0]); code != nil {
				return code.([]byte), nil
			}
			return nil, nil
		}
	}
	code, err := j.callReader1("ReadAccountCode", address)
	if err != nil {
		return nil, err
	}
	if code == nil {
		return nil, nil
	}
	return code.([]byte), nil
}

func (j *Journal) callReader1(method string, args ...any) (any, error) {
	values, err := j.callReader(method, args...)
	if err != nil {
		return nil, err
	}
	switch len(values) {
	case 2:
		if err := resultError(values[1]); err != nil {
			return nil, err
		}
		return resultInterface(values[0]), nil
	case 3:
		if err := resultError(values[2]); err != nil {
			return nil, err
		}
		return resultInterface(values[0]), nil
	default:
		return nil, fmt.Errorf("unexpected %s result count %d", method, len(values))
	}
}

func (j *Journal) callReader(method string, args ...any) ([]reflect.Value, error) {
	m := reflect.ValueOf(j.DB).MethodByName(method)
	if !m.IsValid() {
		return nil, fmt.Errorf("external reader missing %s", method)
	}
	in := make([]reflect.Value, len(args))
	for i, arg := range args {
		in[i] = reflect.ValueOf(arg)
	}
	return m.Call(in), nil
}

func (j *Journal) callReaderWith(method string, arg func(reflect.Type) (reflect.Value, error)) ([]reflect.Value, bool, error) {
	m := reflect.ValueOf(j.DB).MethodByName(method)
	if !m.IsValid() {
		return nil, false, nil
	}
	t := m.Type()
	in := make([]reflect.Value, t.NumIn())
	for i := range in {
		v, err := arg(t.In(i))
		if err != nil {
			return nil, true, err
		}
		in[i] = v
	}
	return m.Call(in), true, nil
}

func addressArg(address types.Address, t reflect.Type) (reflect.Value, error) {
	v := reflect.ValueOf(address)
	if v.Type().AssignableTo(t) {
		return v, nil
	}
	if v.Type().ConvertibleTo(t) {
		return v.Convert(t), nil
	}
	return reflect.Value{}, fmt.Errorf("cannot convert GEVM address to %s", t)
}

func storageKeyArg(index uint256.Int, t reflect.Type) (reflect.Value, error) {
	key := index.Bytes32()
	v := reflect.ValueOf(key)
	if v.Type().AssignableTo(t) {
		return v, nil
	}
	if v.Type().ConvertibleTo(t) {
		return v.Convert(t), nil
	}
	return reflect.Value{}, fmt.Errorf("cannot convert GEVM storage key to %s", t)
}

func accountInfo(raw any, code types.Bytes) (AccountInfo, bool, error) {
	v := reflect.ValueOf(raw)
	if !v.IsValid() {
		return AccountInfo{}, false, nil
	}
	if v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return AccountInfo{}, false, nil
		}
		v = v.Elem()
	}
	balanceField := v.FieldByName("Balance")
	nonceField := v.FieldByName("Nonce")
	incarnationField := v.FieldByName("Incarnation")
	codeHashField := v.FieldByName("CodeHash")
	if !balanceField.IsValid() || !nonceField.IsValid() || !incarnationField.IsValid() || !codeHashField.IsValid() {
		return AccountInfo{}, false, fmt.Errorf("external account type %s is missing required fields", v.Type())
	}
	info := AccountInfo{
		Balance:     balanceField.Interface().(uint256.Int),
		Nonce:       nonceField.Uint(),
		Incarnation: incarnationField.Uint(),
	}
	codeHash, err := hashField(codeHashField)
	if err != nil {
		return AccountInfo{}, false, err
	}
	info.CodeHash = codeHash
	if info.Incarnation == 0 && info.CodeHash != types.B256Zero && info.CodeHash != types.KeccakEmpty {
		info.Incarnation = 1
	}
	if len(code) > 0 && info.CodeHash != types.B256Zero && info.CodeHash != types.KeccakEmpty {
		info.Code = types.Bytes(code)
	}
	return info, true, nil
}

func hashField(v reflect.Value) (types.B256, error) {
	if v.CanInterface() {
		valueMethod := v.MethodByName("Value")
		if valueMethod.IsValid() && valueMethod.Type().NumIn() == 0 && valueMethod.Type().NumOut() == 1 {
			out := valueMethod.Call(nil)[0]
			if out.Type().ConvertibleTo(reflect.TypeOf(types.B256{})) {
				return out.Convert(reflect.TypeOf(types.B256{})).Interface().(types.B256), nil
			}
		}
		if v.Type().ConvertibleTo(reflect.TypeOf(types.B256{})) {
			return v.Convert(reflect.TypeOf(types.B256{})).Interface().(types.B256), nil
		}
	}
	return types.B256Zero, fmt.Errorf("cannot convert external code hash %s", v.Type())
}

func resultInterface(v reflect.Value) any {
	if !v.IsValid() || ((v.Kind() == reflect.Pointer || v.Kind() == reflect.Interface || v.Kind() == reflect.Slice) && v.IsNil()) {
		return nil
	}
	return v.Interface()
}

func resultError(v reflect.Value) error {
	if v.IsNil() {
		return nil
	}
	err, ok := v.Interface().(error)
	if !ok {
		return fmt.Errorf("unexpected error result type %s", v.Type())
	}
	return err
}
