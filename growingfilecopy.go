package main

import (
	"flag"
	"io"
	"log"
	"os"
	"time"
)

// This is the original io.CopyN function improved by using io.CopyBuffer instead of io.Copy
func copyBufferN(dst io.Writer, src io.Reader, b []byte, n int64) (written int64, err error) {
	written, err = io.CopyBuffer(dst, io.LimitReader(src, n), b)
	if written == n {
		return n, nil
	}
	if written < n && err == nil {
		// src stopped early; must have been EOF.
		//err = io.EOF
	}
	return
}

func copy(src string, dst string, chunk int64) {
	// Open the source file in read only mode
	s, err := os.OpenFile(src, os.O_RDONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close() // Close it after we finished the job

	// Remove old destination file to make sure we start from scratch
	err = os.Remove(dst)
	if err != nil {
		log.Println(err)
	}

	// Open destination file in append mode. If file does not exist, create it
	// Mediagrid fsd does not support the syscall o_append. It returns error code 0x10020d
	// See workaround line 49
	d, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY, 0644) //os.O_APPEND|
	if err != nil {
		log.Fatal(err)
	}

	// Start at the beginning of the source file
	var src_offset int64 = 0

	// Run this each second. i increases each run until 5
	// So in fact we wait 6 seconds until we say the file does not grow anymore
	var interval = 6
	i := interval
	for i >= 1 {
		time.Sleep(1000 * time.Millisecond)

		// Get file information from source file
		si, err := s.Stat()
		if err != nil {
			log.Fatal(err)
		}

		// Workaround: Mediagrid does not support syscall o_append
		// Thus we have to set the pointer manually to the eof
		// read file info including size and set pointer
		di, err := d.Stat()
		if err != nil {
			log.Fatal(err)
		}
		_, err = d.Seek(di.Size(), 0)
		if err != nil {
			log.Fatal(err)
		}
		// end of workaround

		log.Println("[check  ] i=", i, "src_offset=", src_offset, "src file size=", si.Size()) // debug output

		// If source file is bigger than our src_offset we copy all the missing data to the end of the destination file
		if si.Size() > src_offset {
			// Set the pointer to the last offset of the source file
			_, err = s.Seek(src_offset, 0)
			if err != nil {
				log.Fatal(err)
			}

			// Copy the data. It returns us how many bytes have been copied
			// Use 2 MB Buffer size
			buf := make([]byte, 2*1024*1024)

			var written int64
			var err error

			if chunk != 0 {
				written, err = copyBufferN(d, s, buf, chunk)
			} else {
				written, err = io.Copy(d, s)
			}
			if err != nil {
				d.Close()
				log.Fatal(err)
			}

			// Set the new offset
			src_offset = src_offset + written

			log.Println("[copy   ] i=", i, "written=", written, "new src_offset=", src_offset) // debug output

			// File has been modified. Restart count countdown
			i = interval

		} else {
			// count down
			i--
		}
	}

	// Close the destination file
	d.Close()
}

func main() {
	log.Println("Started ...")
	var src = flag.String("src", "", "source file path")
	var dst = flag.String("dst", "", "destination file path")
	var chunkSize = flag.Int64("cs", 50*1024*1024, "chunk size in byte")

	flag.Parse()

	copy(*src, *dst, *chunkSize)

	log.Println("... Finished")
}
