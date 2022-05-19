package lib

type Key interface {
	Key() ([]byte, error)
	MustKey() []byte
}

type KeyString struct {
	Value string
}

func (key *KeyString) Key() ([]byte, error) {
	return []byte(key.Value), nil
}

func (key *KeyString) MustKey() []byte {
	return []byte(key.Value)
}

type KeyUint64 struct {
	Value uint64
}

func (key *KeyUint64) Key() ([]byte, error) {
	buf, err := GetBigEndianBytesUint64(key.Value)
	if err != nil {
		return nil, err
	}
	defer ReleaseBytesBuffer(buf)
	return buf.Bytes(), nil
}

func (key *KeyUint64) MustKey() []byte {
	res, err := key.Key()
	if err != nil {
		panic(err)
	}
	return res
}

type KeyInt64 struct {
	Value int64
}

func (key *KeyInt64) Key() ([]byte, error) {
	buf, err := GetBigEndianBytesInt64(key.Value)
	if err != nil {
		return nil, err
	}
	defer ReleaseBytesBuffer(buf)
	return buf.Bytes(), nil
}

func (key *KeyInt64) MustKey() []byte {
	res, err := key.Key()
	if err != nil {
		panic(err)
	}
	return res
}

type KeyByte struct {
	Value byte
}

func (key *KeyByte) Key() ([]byte, error) {
	return []byte{key.Value}, nil
}

func (key *KeyByte) MustKey() []byte {
	return []byte{key.Value}
}

type KeyFloat32 struct {
	Value float32
}

func (key *KeyFloat32) Key() ([]byte, error) {
	buf, err := GetBigEndianBytesFloat32(key.Value)
	if err != nil {
		return nil, err
	}
	defer ReleaseBytesBuffer(buf)
	return buf.Bytes(), nil
}

func (key *KeyFloat32) MustKey() []byte {
	res, err := key.Key()
	if err != nil {
		panic(err)
	}
	return res
}

type KeyBool struct {
	Value bool
}

func (key *KeyBool) Key() ([]byte, error) {
	if key.Value {
		return []byte{1}, nil
	}
	return []byte{0}, nil
}

func (key *KeyBool) MustKey() []byte {
	res, err := key.Key()
	if err != nil {
		panic(err)
	}
	return res
}
