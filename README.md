# idlparser

OMG IDL Parser written in go inspired by [gomme](https://github.com/oleiade/gomme)

## Example

Check `TestParseModule`

```go

func TestParseModule(t *testing.T) {
	tests := []struct {
		input    string
		expected Module
	}{
		{
			input: `module spi {
						bitset idbits {
							bitfield<4> bid; // 4 bits for bus_id
						};

						struct CANFrame {
							@format octet header;
							@format(a=b) idbits id;
						};
					}`,
			expected: Module{
				Name: "spi",
				Content: []ModuleContent{
					bitset.BitSet{
						Name: "idbits",
						Fields: []bitset.Field{
							{
								Name: "bid",
								Type: typeref.BitFieldType{Width: uint8(4)},
							},
						},
					},
					struct_type.Struct{
						Name: "CANFrame",
						Fields: []struct_type.Field{
							{
								Name:        "header",
								Annotations: []annotation.Annotation{{Name: "format"}},
								Type:        typeref.OctetType{},
							},
							{
								Name:        "id",
								Annotations: []annotation.Annotation{{Name: "format", Values: map[string]string{"a": "b"}}},
								Type:        typeref.TypeName{Name: "idbits"},
							},
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		result := Parse(test.input)
		require.Equal(t, test.expected, result.Output)
	}
}

```