// Export type instead of the underlying gRPC struct
// May need to update struct in this file if the protobuf struct is being changed

package gtoolkits

const (
	ADDRESS = "localhost:50052"
	TIMEOUT = 2
)

type TFRecord struct {
	Keyword string
	Weight  float64
}
