package parser_test

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/m-lab/etl/parser"
	"github.com/m-lab/etl/schema"
)

/*
// TODO: IPv6 tests
func TestParseFirstLine(t *testing.T) {
	protocol, dest_ip, server_ip, err := parser.ParseFirstLine("traceroute [(64.86.132.76:33461) -> (98.162.212.214:53849)], protocol icmp, algo exhaustive, duration 19 s")
	if dest_ip != "98.162.212.214" || server_ip != "64.86.132.76" || protocol != "icmp" || err != nil {
		t.Errorf("Error in parsing the first line!\n")
		return
	}

	protocol, dest_ip, server_ip, err = parser.ParseFirstLine("Exception : [ERROR](Probe.cc, 109)Can't send the probe : Invalid argument")
	if err == nil {
		t.Errorf("Error in parsing the first line!\n")
		return
	}

}

func TestCreateTestId(t *testing.T) {
	test_id := parser.CreateTestId("20170501T000000Z-mlab1-acc02-paris-traceroute-0000.tgz", "20170501T23:53:10Z-98.162.212.214-53849-64.86.132.75-42677.paris")
	if test_id != "2017/05/01/mlab1.acc02/20170501T23:53:10Z-98.162.212.214-53849-64.86.132.75-42677.paris.gz" {
		fmt.Println(test_id)
		t.Errorf("Error in creating test id!\n")
		return
	}
}

func TestParseLegacyFormatData(t *testing.T) {
	rawData, err := ioutil.ReadFile("testdata/20160112T00:45:44Z_ALL27409.paris")
	if err != nil {
		fmt.Println("cannot load test data")
		return
	}
	cashedTest, err := parser.Parse(nil, "testdata/20160112T00:45:44Z_ALL27409.paris", "", rawData, "pt-daily")
	if err != nil {
		t.Fatalf(err.Error())
	}
	if len(cashedTest.Hops) != 9 {
		t.Fatalf("Do not process hops correctly.")
	}
	if cashedTest.LogTime.Unix() != 1452559544 {
		t.Fatalf("Do not process log time correctly.")
	}
	if cashedTest.LastValidHopLine != "ExpectedDestIP" {
		t.Fatalf("Did not reach expected destination.")
	}
}
*/
func TestPTParser(t *testing.T) {
	rawData, err := ioutil.ReadFile("testdata/20170320T23:53:10Z-172.17.94.34-33456-74.125.224.100-33457.paris")
	cashedTest, err := parser.Parse(nil, "testdata/20170320T23:53:10Z-172.17.94.34-33456-74.125.224.100-33457.paris", "", rawData, "pt-daily")
	if err != nil {
		t.Fatalf(err.Error())
	}
	if cashedTest.LogTime.Unix() != 1490053990 {
		t.Fatalf("Do not process log time correctly.")
	}

	if cashedTest.Source.Ip != "172.17.94.34" {
		t.Fatalf("Wrong results for Server IP.")
	}

	if cashedTest.Destination.Ip != "74.125.224.100" {
		t.Fatalf("Wrong results for Client IP.")
	}

	log.Printf("%+v", cashedTest.Hops)

	// TODO(dev): reformat these individual values to be more readable.
	expected_hops := []schema.ScamperHop{
		schema.ScamperHop{Source:{Ip:"64.233.174.109", City:"", CountryCode:"", Hostname:"sr05-te1-8.nuq04.net.google.com"}, Linkc:0, Links:[{HopDstIp:"74.125.224.100", TTL:0, Probes:[{Flowid:0, Rtt:[]float64{0.895}}]}]},
		schema.ScamperHop{Source:{Ip:72.14.232.136 City: CountryCode: Hostname:bb01-ae7.nuq04.net.google.com} Linkc:0 Links:[{HopDstIp:64.233.174.109 TTL:0 Probes:[{Flowid:0 Rtt:[1.614]}]}]},
		schema.ScamperHop{Source:{Ip:72.14.232.136 City: CountryCode: Hostname:bb01-ae7.nuq04.net.google.com} Linkc:0 Links:[{HopDstIp:64.233.174.109 TTL:0 Probes:[{Flowid:0 Rtt:[1.614]}]}]},
		schema.ScamperHop{Source:{Ip:72.14.232.136 City: CountryCode: Hostname:bb01-ae7.nuq04.net.google.com} Linkc:0 Links:[{HopDstIp:64.233.174.109 TTL:0 Probes:[{Flowid:0 Rtt:[1.614]}]}]},
		schema.ScamperHop{Source:{Ip:72.14.232.136 City: CountryCode: Hostname:bb01-ae7.nuq04.net.google.com} Linkc:0 Links:[{HopDstIp:64.233.174.109 TTL:0 Probes:[{Flowid:0 Rtt:[1.614]}]}]},
		schema.ScamperHop{Source:{Ip:216.239.49.250 City: CountryCode: Hostname:bb01-ae3.nuq04.net.google.com} Linkc:0 Links:[{HopDstIp:64.233.174.109 TTL:0 Probes:[{Flowid:0 Rtt:[1.614]}]}]},
		schema.ScamperHop{Source:{Ip:216.239.49.250 City: CountryCode: Hostname:bb01-ae3.nuq04.net.google.com} Linkc:0 Links:[{HopDstIp:64.233.174.109 TTL:0 Probes:[{Flowid:0 Rtt:[1.614]}]}]},
		schema.ScamperHop{Source:{Ip:216.239.49.250 City: CountryCode: Hostname:bb01-ae3.nuq04.net.google.com} Linkc:0 Links:[{HopDstIp:64.233.174.109 TTL:0 Probes:[{Flowid:0 Rtt:[1.614]}]}]},
		schema.ScamperHop{Source:{Ip:216.239.49.250 City: CountryCode: Hostname:bb01-ae3.nuq04.net.google.com} Linkc:0 Links:[{HopDstIp:64.233.174.109 TTL:0 Probes:[{Flowid:0 Rtt:[1.614]}]}]},
		schema.ScamperHop{Source:{Ip:216.239.49.250 City: CountryCode: Hostname:bb01-ae3.nuq04.net.google.com} Linkc:0 Links:[{HopDstIp:64.233.174.109 TTL:0 Probes:[{Flowid:0 Rtt:[1.614]}]}]},
		schema.ScamperHop{Source:{Ip:216.239.49.250 City: CountryCode: Hostname:bb01-ae3.nuq04.net.google.com} Linkc:0 Links:[{HopDstIp:64.233.174.109 TTL:0 Probes:[{Flowid:0 Rtt:[1.614]}]}]},
		schema.ScamperHop{Source:{Ip:216.239.49.250 City: CountryCode: Hostname:bb01-ae3.nuq04.net.google.com} Linkc:0 Links:[{HopDstIp:64.233.174.109 TTL:0 Probes:[{Flowid:0 Rtt:[1.614]}]}]},
		schema.ScamperHop{Source:{Ip:72.14.196.8 City: CountryCode: Hostname:pr02-xe-3-0-1.pao03.net.google.com} Linkc:0 Links:[{HopDstIp:72.14.232.136 TTL:0 Probes:[{Flowid:0 Rtt:[1.693]}]}]},
		schema.ScamperHop{Source:{Ip:72.14.196.8 City: CountryCode: Hostname:pr02-xe-3-0-1.pao03.net.google.com} Linkc:0 Links:[{HopDstIp:72.14.232.136 TTL:0 Probes:[{Flowid:0 Rtt:[1.693]}]}]},
		schema.ScamperHop{Source:{Ip:72.14.196.8 City: CountryCode: Hostname:pr02-xe-3-0-1.pao03.net.google.com} Linkc:0 Links:[{HopDstIp:72.14.232.136 TTL:0 Probes:[{Flowid:0 Rtt:[1.693]}]}]},
		schema.ScamperHop{Source:{Ip:72.14.196.8 City: CountryCode: Hostname:pr02-xe-3-0-1.pao03.net.google.com} Linkc:0 Links:[{HopDstIp:72.14.232.136 TTL:0 Probes:[{Flowid:0 Rtt:[1.693]}]}]},
		schema.ScamperHop{Source:{Ip:72.14.218.190 City: CountryCode: Hostname:pr01-xe-7-1-0.pao03.net.google.com} Linkc:0 Links:[{HopDstIp:216.239.49.250 TTL:0 Probes:[{Flowid:0 Rtt:[1.386]}]}]},
		schema.ScamperHop{Source:{Ip:72.14.218.190 City: CountryCode: Hostname:pr01-xe-7-1-0.pao03.net.google.com} Linkc:0 Links:[{HopDstIp:216.239.49.250 TTL:0 Probes:[{Flowid:0 Rtt:[1.386]}]}]},
		schema.ScamperHop{Source:{Ip:72.14.218.190 City: CountryCode: Hostname:pr01-xe-7-1-0.pao03.net.google.com} Linkc:0 Links:[{HopDstIp:216.239.49.250 TTL:0 Probes:[{Flowid:0 Rtt:[1.386]}]}]},
		schema.ScamperHop{Source:{Ip:72.14.218.190 City: CountryCode: Hostname:pr01-xe-7-1-0.pao03.net.google.com} Linkc:0 Links:[{HopDstIp:216.239.49.250 TTL:0 Probes:[{Flowid:0 Rtt:[1.386]}]}]},
		schema.ScamperHop{Source:{Ip:72.14.218.190 City: CountryCode: Hostname:pr01-xe-7-1-0.pao03.net.google.com} Linkc:0 Links:[{HopDstIp:216.239.49.250 TTL:0 Probes:[{Flowid:0 Rtt:[1.386]}]}]},
		schema.ScamperHop{Source:{Ip:72.14.218.190 City: CountryCode: Hostname:pr01-xe-7-1-0.pao03.net.google.com} Linkc:0 Links:[{HopDstIp:216.239.49.250 TTL:0 Probes:[{Flowid:0 Rtt:[1.386]}]}]},
		schema.ScamperHop{Source:{Ip:72.14.218.190 City: CountryCode: Hostname:pr01-xe-7-1-0.pao03.net.google.com} Linkc:0 Links:[{HopDstIp:216.239.49.250 TTL:0 Probes:[{Flowid:0 Rtt:[1.386]}]}]},
		schema.ScamperHop{Source:{Ip:172.25.253.46 City: CountryCode: Hostname:us-mtv-ply1-br1-xe-1-1-0-706.n.corp.google.com} Linkc:0 Links:[{HopDstIp:72.14.196.8 TTL:0 Probes:[{Flowid:0 Rtt:[0.556]}]}]},
		schema.ScamperHop{Source:{Ip:172.25.253.46 City: CountryCode: Hostname:us-mtv-ply1-br1-xe-1-1-0-706.n.corp.google.com} Linkc:0 Links:[{HopDstIp:72.14.196.8 TTL:0 Probes:[{Flowid:0 Rtt:[0.556]}]}]},
		schema.ScamperHop{Source:{Ip:172.25.253.46 City: CountryCode: Hostname:us-mtv-ply1-br1-xe-1-1-0-706.n.corp.google.com} Linkc:0 Links:[{HopDstIp:72.14.196.8 TTL:0 Probes:[{Flowid:0 Rtt:[0.556]}]}]},
		schema.ScamperHop{Source:{Ip:172.25.253.46 City: CountryCode: Hostname:us-mtv-ply1-br1-xe-1-1-0-706.n.corp.google.com} Linkc:0 Links:[{HopDstIp:72.14.196.8 TTL:0 Probes:[{Flowid:0 Rtt:[0.556]}]}]},
		schema.ScamperHop{Source:{Ip:172.25.253.46 City: CountryCode: Hostname:us-mtv-ply1-br1-xe-1-1-0-706.n.corp.google.com} Linkc:0 Links:[{HopDstIp:72.14.218.190 TTL:0 Probes:[{Flowid:0 Rtt:[0.53]}]}]},
		schema.ScamperHop{Source:{Ip:172.25.253.46 City: CountryCode: Hostname:us-mtv-ply1-br1-xe-1-1-0-706.n.corp.google.com} Linkc:0 Links:[{HopDstIp:72.14.218.190 TTL:0 Probes:[{Flowid:0 Rtt:[0.53]}]}]},
		schema.ScamperHop{Source:{Ip:172.25.253.46 City: CountryCode: Hostname:us-mtv-ply1-br1-xe-1-1-0-706.n.corp.google.com} Linkc:0 Links:[{HopDstIp:72.14.218.190 TTL:0 Probes:[{Flowid:0 Rtt:[0.53]}]}]},
		schema.ScamperHop{Source:{Ip:172.25.253.46 City: CountryCode: Hostname:us-mtv-ply1-br1-xe-1-1-0-706.n.corp.google.com} Linkc:0 Links:[{HopDstIp:72.14.218.190 TTL:0 Probes:[{Flowid:0 Rtt:[0.53]}]}]},
		schema.ScamperHop{Source:{Ip:172.25.253.46 City: CountryCode: Hostname:us-mtv-ply1-br1-xe-1-1-0-706.n.corp.google.com} Linkc:0 Links:[{HopDstIp:72.14.218.190 TTL:0 Probes:[{Flowid:0 Rtt:[0.53]}]}]},
		schema.ScamperHop{Source:{Ip:172.25.253.46 City: CountryCode: Hostname:us-mtv-ply1-br1-xe-1-1-0-706.n.corp.google.com} Linkc:0 Links:[{HopDstIp:72.14.218.190 TTL:0 Probes:[{Flowid:0 Rtt:[0.53]}]}]},
		schema.ScamperHop{Source:{Ip:172.25.253.46 City: CountryCode: Hostname:us-mtv-ply1-br1-xe-1-1-0-706.n.corp.google.com} Linkc:0 Links:[{HopDstIp:72.14.218.190 TTL:0 Probes:[{Flowid:0 Rtt:[0.53]}]}]},
		schema.ScamperHop{Source:{Ip:172.25.252.166 City: CountryCode: Hostname:us-mtv-ply1-bb1-tengigabitethernet2-3.n.corp.google.com} Linkc:0 Links:[{HopDstIp:172.25.253.46 TTL:0 Probes:[{Flowid:0 Rtt:[0.343]}]}]} ,
		schema.ScamperHop{Source:{Ip:172.25.252.172 City: CountryCode: Hostname:us-mtv-cl4-core1-gigabitethernet1-1.n.corp.google.com} Linkc:0 Links:[{HopDstIp:172.25.252.166 TTL:0 Probes:[{Flowid:0 Rtt:[0.501]}]}]},
		schema.ScamperHop{Source:{Ip:172.17.95.252 City: CountryCode: Hostname:172.17.95.252} Linkc:0 Links:[{HopDstIp:172.25.252.172 TTL:0 Probes:[{Flowid:0 Rtt:[0.407]}]}]},
		schema.ScamperHop{Source:{Ip:172.17.94.34 City: CountryCode: Hostname:} Linkc:0 Links:[{HopDstIp:172.17.95.252 TTL:0 Probes:[{Flowid:0 Rtt:[0.376]}]}]},
	}
	if len(cashedTest.Hops) != len(expected_hops) {
		t.Fatalf("Wrong results for PT hops!")
	}
	/*
		for i := 0; i < len(cashedTest.Hops); i++ {
			if !reflect.DeepEqual(cashedTest.Hops[i], expected_hops[i]) {
				fmt.Println(i)
				fmt.Printf("Here is expected    : %v\n", expected_hops[i])
				fmt.Printf("Here is what is real: %v\n", cashedTest.Hops[i])
				t.Fatalf("Wrong results for PT hops!")
			}
		}*/
}

/*
func TestPTInserter(t *testing.T) {
	ins := &inMemoryInserter{}
	pt := parser.NewPTParser(ins)
	rawData, err := ioutil.ReadFile("testdata/20170320T23:53:10Z-172.17.94.34-33456-74.125.224.100-33457.paris")
	if err != nil {
		t.Fatalf("cannot read testdata.")
	}
	meta := map[string]bigquery.Value{"filename": "gs://fake-bucket/fake-archive.tgz"}
	err = pt.ParseAndInsert(meta, "testdata/20170320T23:53:10Z-172.17.94.34-33456-74.125.224.100-33457.paris", rawData)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if ins.RowsInBuffer() != 38 {
		fmt.Println(ins.RowsInBuffer())
		t.Fatalf("Number of rows in PT table is wrong.")
	}
	fmt.Println(ins.data[0])

	expectedValues := &schema.PT{}

	// Copy ParseTime from actual output before using DeepEqual.
	expectedValues.ParseTime = (ins.data[0].(schema.PT)).ParseTime
	if !reflect.DeepEqual(ins.data[0], *expectedValues) {
		fmt.Printf("Here is expected    : %v\n", expectedValues)
		fmt.Printf("Here is what is real: %v\n", ins.data[0])
		t.Errorf("Not the expected values:")
	}
}

func TestPTPollutionCheck(t *testing.T) {
	ins := &inMemoryInserter{}
	pt := parser.NewPTParser(ins)

	tests := []struct {
		fileName             string
		expectedBufferedTest int
		expectedNumRows      int
	}{
		{
			fileName:             "testdata/PT/20171208T00:00:04Z-35.188.101.1-40784-173.205.3.38-9090.paris",
			expectedBufferedTest: 1,
			expectedNumRows:      0,
		},
		{
			fileName: "testdata/PT/20171208T00:00:04Z-37.220.21.130-5667-173.205.3.43-42487.paris",
			// The second test reached expected destIP, and was inserted into BigQuery table.
			// The buffer has only the first test.
			expectedBufferedTest: 1,
			expectedNumRows:      16,
		},
		{
			fileName: "testdata/PT/20171208T00:00:14Z-139.60.160.135-2023-173.205.3.44-1101.paris",
			// The first test was detected that it was polluted by the third test.
			// expectedBufferedTest is 0, which means pollution detected and test removed.
			expectedBufferedTest: 0,
			// The third test reached its destIP and was inserted into BigQuery.
			expectedNumRows: 29,
		},
		{
			fileName: "testdata/PT/20171208T00:00:14Z-76.227.226.149-37156-173.205.3.37-52156.paris",
			// The 4th test was buffered.
			expectedBufferedTest: 1,
			expectedNumRows:      29,
		},
		{
			fileName: "testdata/PT/20171208T22:03:54Z-104.198.139.160-60574-163.22.28.37-7999.paris",
			// The 5th test was buffered too.
			expectedBufferedTest: 2,
			expectedNumRows:      29,
		},
		{
			fileName: "testdata/PT/20171208T22:03:59Z-139.60.160.135-1519-163.22.28.44-1101.paris",
			// The 5th test was detected that was polluted by the 6th test.
			// It was removed from buffer (expectedBufferedTest drop from 2 to 1).
			// Buffer contains the 4th test now.
			expectedBufferedTest: 1,
			// The 6th test reached its destIP and was inserted into BigQuery.
			expectedNumRows: 46,
		},
	}

	// Process the tests
	for _, test := range tests {
		rawData, err := ioutil.ReadFile(test.fileName)
		if err != nil {
			t.Fatalf("cannot read testdata.")
		}
		err = pt.ParseAndInsert(nil, test.fileName, rawData)
		if err != nil {
			t.Fatalf(err.Error())
		}
		if pt.NumBufferedTests() != test.expectedBufferedTest {
			t.Fatalf("Data not buffered correctly")
		}
		if ins.RowsInBuffer() != test.expectedNumRows {
			t.Fatalf("Data not inserted into BigQuery correctly.")
		}
	}

	// Insert the 4th test in the buffer to BigQuery.
	pt.ProcessLastTests()
	if ins.RowsInBuffer() != 56 {
		t.Fatalf("Data not inserted into BigQuery correctly.")
	}

}

func TestPTEmptyTest(t *testing.T) {
	rawData, err := ioutil.ReadFile("testdata/20180201T07:57:37Z-125.212.217.215-56622-208.177.76.115-9100.paris")
	if err != nil {
		fmt.Println("cannot load test data")
		return
	}
	_, parseErr := parser.Parse(nil, "testdata/20180201T07:57:37Z-125.212.217.215-56622-208.177.76.115-9100.paris", "", rawData, "pt-daily")
	if parseErr.Error() != "Empty Test." {
		t.Fatalf("Do not handle empty test correctly.")
	}
}
*/
