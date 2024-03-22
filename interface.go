package interfaceio

import (
	"fmt"
	"io"
	"os"
)

// Definition of custom Writer
type CustomWriter struct {
	w io.Writer /*
		That's partially correct! The io.Writer interface provides a foundation for writing data, but the exact output format depends on the specific type of io.Writer you use. Here's a breakdown:

		Flexibility Within Constraints:  While io.Writer allows you to output various data streams, the final format is ultimately determined by the underlying writer you pass to your CustomWriter.

		Examples of io.Writer Types:

		A file (*os.File) writes bytes to a file on disk.
		Network connections (net.Conn) send data packets over a network.
		bytes.Buffer acts as an in-memory buffer, accumulating the data you write.
		You could even create your own custom io.Writer that performs specific formatting.
		CustomWriter's Role:  Your CustomWriter doesn't directly control the final output format. What it offers is:

		The ability to:
		Change Output Destination: You can redirect the data to different locations (files, network, etc.) based on the io.Writer you provide.
		Potentially Modify Data: Within the Write method of your CustomWriter, you can potentially transform the data (e.g., encryption, adding headers) before it goes to the final writer
	*/
}

// Implementation of the NewWriter function

func NewWriter(w io.Writer) *CustomWriter {
	return &CustomWriter{
		w: w,
	}
}

// Implementation of the Write function
func (cw *CustomWriter) Write(bs []byte) (int, error) {
	return cw.w.Write(bs)
}

// Definition of custom Reader
type CustomReader struct {
	file         *os.File // pointer to our file so we can access the content
	currentIndex int64    // Current position in file
	bufferSize   int      // bufferSize signifies how much data we can read at once
}

// This serves as a construction for the Reader
func NewCustomReader(bufferSize int) (*CustomReader, error) {
	if len(os.Args) < 2 {
		fmt.Println("You did not specify a file")
		os.Exit(1)
	}
	file, err := os.Open(os.Args[1])
	if err != nil {
		return nil, err
	}

	return &CustomReader{
		file:         file,
		currentIndex: 0,
		bufferSize:   bufferSize,
	}, nil
}

// Definition of Read function to implement Reader interface
// We need to adjust the logic of the Read function to fit our needs
// We need to use Open function to read the contents of our file
// This was way harder than expected but I now understand how Read function within the os package works
func (cw *CustomReader) Read(bs []byte) (int, error) {
	// Handle if you've reached the end of the file
	fileInfo, err := cw.file.Stat()
	if err != nil {
		fmt.Println("Error occurred with file")
		os.Exit(1)
	}
	if cw.currentIndex >= fileInfo.Size() {
		return 0, io.EOF
	}

	// Determine bytes to read as long as the buffer can store the data
	bytesToRead := cw.bufferSize
	if bytesToRead > len(bs) {
		bytesToRead = len(bs)
	}

	// Read from where we left off
	n, err := cw.file.ReadAt(bs[:bytesToRead], cw.currentIndex)
	if err != nil && err != io.EOF {
		return n, err // Return read bytes even if there's a non-EOF error
	}

	// Update position, return bytes read and any error
	cw.currentIndex += int64(n)
	return n, err
}

// Function to close file since we have a * to os.File
func (cw *CustomReader) CloseFile() error {
	return cw.file.Close()
}
