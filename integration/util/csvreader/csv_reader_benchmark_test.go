package csvreader

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"flag"
	"io"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/golang/snappy"
)

var testFlags struct {
	csvFile                  string
	readBufferBytes          int
	cachePerChanTransmission int
	chanCacheSize            int
	snappyCompressed         bool
}

func init() {
	flag.StringVar(&testFlags.csvFile, "csv-file", "", "Input csv file for testing.")
	flag.IntVar(&testFlags.readBufferBytes, "readbuf", 0, "Create specified size buffer for reading if > 0, otherwise use Reader's default.")
	flag.IntVar(&testFlags.cachePerChanTransmission, "cache-per-trans", 500, "If cache before chan transmission, how many caches per trans.")
	flag.IntVar(&testFlags.chanCacheSize, "cached-chan", 100, "Chan cache size. 0 if blocked chan.")
	flag.BoolVar(&testFlags.snappyCompressed, "snappy-compressed", false, "Whether the csv snappy compressed or not.")
}

func makeStringChan() chan string {
	return make(chan string, testFlags.chanCacheSize)
}

func makeStringSliceChan() chan []string {
	return make(chan []string, testFlags.chanCacheSize)
}

func makeRecordChan() chan []string {
	return makeStringSliceChan()
}

func makeRecordSliceChan() chan [][]string {
	return make(chan [][]string, testFlags.chanCacheSize)
}

func stringConsumer(in <-chan string) {
	for {
		_, ok := <-in
		if !ok {
			break
		}
	}
}

func stringSliceConsumer(in <-chan []string) {
	for {
		_, ok := <-in
		if !ok {
			break
		}
	}
}

func recordConsumer(in <-chan []string) {
	stringSliceConsumer(in)
}

func recordSliceConsumer(in <-chan [][]string) {
	for {
		_, ok := <-in
		if !ok {
			break
		}
	}
}

func makeCompressedReader(r io.Reader) io.Reader {
	if testFlags.snappyCompressed {
		return snappy.NewReader(r)
	}
	return r
}

func makeBufferedReader(r io.Reader) io.Reader {
	if testFlags.readBufferBytes <= 0 {
		return r
	}
	return bufio.NewReaderSize(r, testFlags.readBufferBytes)
}

func makeWrappedReader(r io.Reader) io.Reader {
	return makeCompressedReader(makeBufferedReader(r))
}

func BenchmarkConsumingBytesInPlaceFromBufioScan(b *testing.B) {
	for i := 0; i < b.N; i++ {

		f, err := os.Open(testFlags.csvFile)
		defer f.Close()
		if err != nil {
			b.Error(err)
		}

		scanner := bufio.NewScanner(makeWrappedReader(f))
		for scanner.Scan() {
			_ = scanner.Bytes()
		}
	}
}

func BenchmarkConsumingTextInPlaceFromBufioScan(b *testing.B) {
	for i := 0; i < b.N; i++ {

		f, err := os.Open(testFlags.csvFile)
		defer f.Close()
		if err != nil {
			b.Error(err)
		}

		scanner := bufio.NewScanner(makeWrappedReader(f))
		for scanner.Scan() {
			_ = scanner.Text()
		}
	}
}

func BenchmarkConsumingTextFromBufioScan(b *testing.B) {
	for i := 0; i < b.N; i++ {

		f, err := os.Open(testFlags.csvFile)
		defer f.Close()
		if err != nil {
			b.Error(err)
		}

		stringChan := makeStringChan()
		go stringConsumer(stringChan)

		scanner := bufio.NewScanner(makeWrappedReader(f))
		for scanner.Scan() {
			stringChan <- scanner.Text()
		}
		close(stringChan)
	}
}

func BenchmarkConsumingTextSliceFromBufioScan(b *testing.B) {
	cacheCount := testFlags.cachePerChanTransmission

	for i := 0; i < b.N; i++ {

		f, err := os.Open(testFlags.csvFile)
		defer f.Close()
		if err != nil {
			b.Error(err)
		}

		stringSliceChan := makeStringSliceChan()
		go stringSliceConsumer(stringSliceChan)

		scanner := bufio.NewScanner(makeWrappedReader(f))

		stringSliceCache := make([]string, 0, cacheCount)
		for scanner.Scan() {
			stringSliceCache = append(stringSliceCache, scanner.Text())
			if len(stringSliceCache) >= cacheCount {
				stringSliceChan <- stringSliceCache
				stringSliceCache = make([]string, 0, cacheCount)
			}
		}
		close(stringSliceChan)
	}
}

func BenchmarkConsumingTextSliceWithPoolFromBufioScan(b *testing.B) {
	cacheCount := testFlags.cachePerChanTransmission

	for i := 0; i < b.N; i++ {

		f, err := os.Open(testFlags.csvFile)
		defer f.Close()
		if err != nil {
			b.Error(err)
		}

		var bufPool = sync.Pool{
			New: func() interface{} {
				// The Pool's New function should generally only return pointer
				// types, since a pointer can be put into the return interface
				// value without an allocation:
				s := make([]string, 0, cacheCount)
				return &s
			},
		}

		stringSliceChan := make(chan *[]string, testFlags.chanCacheSize)
		go func(in <-chan *[]string) {
			for {
				s, ok := <-in
				if !ok {
					break
				}
				*s = (*s)[:0]
				bufPool.Put(s)
			}
		}(stringSliceChan)

		scanner := bufio.NewScanner(makeWrappedReader(f))

		stringSliceCache := bufPool.Get().(*[]string)
		for scanner.Scan() {
			*stringSliceCache = append(*stringSliceCache, scanner.Text())
			if len(*stringSliceCache) >= cacheCount {
				stringSliceChan <- stringSliceCache
				stringSliceCache = bufPool.Get().(*[]string)
			}
		}
		close(stringSliceChan)
	}
}

func BenchmarkConsumingRecordFromBufioScanAndBytesSplit(b *testing.B) {
	for i := 0; i < b.N; i++ {

		f, err := os.Open(testFlags.csvFile)
		defer f.Close()
		if err != nil {
			b.Error(err)
		}

		recordChan := makeRecordChan()
		go recordConsumer(recordChan)

		scanner := bufio.NewScanner(makeWrappedReader(f))
		for scanner.Scan() {
			bytes2DArray := bytes.Split(scanner.Bytes(), []byte{','})
			record := make([]string, len(bytes2DArray))
			for j := 0; j < len(bytes2DArray); j++ {
				record[j] = string(bytes2DArray[j])
			}
			recordChan <- record
		}
		close(recordChan)
	}
}
func BenchmarkConsumingRecordFromBufioScanAndSplit(b *testing.B) {
	for i := 0; i < b.N; i++ {

		f, err := os.Open(testFlags.csvFile)
		defer f.Close()
		if err != nil {
			b.Error(err)
		}

		recordChan := makeRecordChan()
		go recordConsumer(recordChan)

		scanner := bufio.NewScanner(makeWrappedReader(f))
		for scanner.Scan() {
			s := scanner.Text()
			recordChan <- strings.Split(s, ",")
		}
		close(recordChan)
	}
}

func benchmarkCSVPkg(b *testing.B, csvFile string, reuseRecord bool) {
	for i := 0; i < b.N; i++ {

		f, err := os.Open(csvFile)
		defer f.Close()
		if err != nil {
			b.Error(err)
		}

		r := csv.NewReader(makeWrappedReader(f))
		r.ReuseRecord = reuseRecord
		r.FieldsPerRecord = -1 // disable fields count check

		recordChan := makeRecordChan()
		go recordConsumer(recordChan)

		for {
			record, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				b.Error(err)
			}
			recordChan <- record
		}
		close(recordChan)
	}

}

func BenchmarkConsumingRecordFromCSVRead(b *testing.B) {
	benchmarkCSVPkg(b, testFlags.csvFile, false)
}

func BenchmarkConsumingRecordFromCSVReadReuseRecord(b *testing.B) {
	benchmarkCSVPkg(b, testFlags.csvFile, true)
}

func BenchmarkConsumingRecordSliceFromBufioScanAndSplit(b *testing.B) {
	cacheCount := testFlags.cachePerChanTransmission

	for i := 0; i < b.N; i++ {

		f, err := os.Open(testFlags.csvFile)
		defer f.Close()
		if err != nil {
			b.Error(err)
		}

		recordSliceChan := makeRecordSliceChan()
		go recordSliceConsumer(recordSliceChan)

		scanner := bufio.NewScanner(makeWrappedReader(f))

		recordSliceCache := make([][]string, 0, cacheCount)
		for scanner.Scan() {
			s := scanner.Text()
			recordSliceCache = append(recordSliceCache, strings.Split(s, ","))
			if len(recordSliceCache) >= cacheCount {
				recordSliceChan <- recordSliceCache
				recordSliceCache = make([][]string, 0, cacheCount)
			}
		}
		close(recordSliceChan)
	}
}

func BenchmarkConsumingRecordSliceFromCSVReadReuseRecord(b *testing.B) {
	cacheCount := testFlags.cachePerChanTransmission

	for i := 0; i < b.N; i++ {

		f, err := os.Open(testFlags.csvFile)
		defer f.Close()
		if err != nil {
			b.Error(err)
		}

		r := csv.NewReader(makeWrappedReader(f))
		r.ReuseRecord = true
		r.FieldsPerRecord = -1 // disable fields count check

		recordSliceChan := makeRecordSliceChan()
		go recordSliceConsumer(recordSliceChan)

		recordSliceCache := make([][]string, 0, cacheCount)
		for {
			record, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				b.Error(err)
			}
			recordSliceCache = append(recordSliceCache, record)
			if len(recordSliceCache) >= cacheCount {
				recordSliceChan <- recordSliceCache
				recordSliceCache = make([][]string, 0, cacheCount)
			}
		}
		close(recordSliceChan)
	}

}

func BenchmarkGenerateRecordByConsumingTextSliceFromBufioScan(b *testing.B) {
	cacheCount := testFlags.cachePerChanTransmission
	parallelTransformCount := 3

	for i := 0; i < b.N; i++ {

		f, err := os.Open(testFlags.csvFile)
		defer f.Close()
		if err != nil {
			b.Error(err)
		}

		stringSliceChan := makeStringSliceChan()

		wg := sync.WaitGroup{}
		for j := 0; j < parallelTransformCount; j++ {
			wg.Add(1)
			go func(in <-chan []string) {
				for {
					ss, ok := <-in
					if !ok {
						break
					}

					for _, s := range ss {
						_ = strings.Split(s, ",")
					}
				}
				wg.Done()
			}(stringSliceChan)
		}

		scanner := bufio.NewScanner(makeWrappedReader(f))

		stringSliceCache := make([]string, 0, cacheCount)
		for scanner.Scan() {
			stringSliceCache = append(stringSliceCache, scanner.Text())
			if len(stringSliceCache) >= cacheCount {
				stringSliceChan <- stringSliceCache
				stringSliceCache = make([]string, 0, cacheCount)
			}
		}
		close(stringSliceChan)

		wg.Wait()
	}
}

func BenchmarkGenerateRecordByConsumingTextSliceFromLinesAsyncReader(b *testing.B) {

	parallelTransformCount := 3
	options := Options{
		MaxCacheCount: testFlags.chanCacheSize,
		MinReadCount:  testFlags.cachePerChanTransmission,
	}
	if testFlags.snappyCompressed {
		options.Compression = CompressionTypeSnappy
	}

	for i := 0; i < b.N; i++ {

		l := NewLinesAsyncReader(testFlags.csvFile, &options)

		wg := sync.WaitGroup{}
		for j := 0; j < parallelTransformCount; j++ {
			wg.Add(1)
			go func() {
				for {
					ss, ok := l.ReadLines()
					if !ok {
						break
					}

					for _, s := range ss {
						_ = strings.Split(s, ",")
					}
				}
				wg.Done()
			}()
		}

		l.Start()
		wg.Wait()

		if err := l.Err(); err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkConsumingRecordSliceFromRecordsAsyncReader(b *testing.B) {

	options := Options{
		MaxCacheCount: testFlags.chanCacheSize,
		MinReadCount:  testFlags.cachePerChanTransmission,
	}
	if testFlags.snappyCompressed {
		options.Compression = CompressionTypeSnappy
	}

	for i := 0; i < b.N; i++ {

		r := NewRecordsAsyncReader(testFlags.csvFile, &options)

		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			for {
				_, ok := r.ReadRecords()
				if !ok {
					break
				}
			}
			wg.Done()
		}()

		r.Start()
		wg.Wait()

		if err := r.Err(); err != nil {
			b.Error(err)
		}
	}
}

// Uncomment below TestMain if want to see glog output in `go test` process.
// func TestMain(m *testing.M) {
// 	flag.Set("alsologtostderr", "true")
// 	flag.Set("log_dir", ".")
// 	flag.Set("v", "2")
// 	flag.Parse()

// 	ret := m.Run()
// 	os.Exit(ret)
// }
