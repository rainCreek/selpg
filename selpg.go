package main

import (
  "fmt"
  "os"
  "io"
  "strconv"
  "os/exec"
  "bufio"
)


/*================================= types =========================*/

type sp_args struct {
  start_page int
  end_page int
  in_filename string
  page_len int   /* default value, can be overriden by "-l number" on command line */
  page_type int   /* 'l' for lines-delimited, 'f' for form-feed-delimited */
  print_dest string
}

const INBUFSIZ = 16 * 1024

/*================================= globals =======================*/

var progname string  /* program name, for error messages */

/*================================= main()=== =====================*/

func main() {
  sa := sp_args { 
    start_page:-1, 
    end_page:-1, 
    page_len:72, 
    page_type:1, 
    in_filename: "", 
    print_dest:"", 
  } 

  progname = os.Args[0] 
  process_args(&sa) 
  process_input(sa) 

  //return 0;
}


/*================================= process_args() ================*/

func process_args(psa *sp_args) {
  var s1 string
  //var s2 string
  var argno int
  var i int

  /* check the command-line arguments for validity */
  if len(os.Args) < 3 {   /* Not enough args, minimum command is "selpg -sstartpage -eend_page"  */
    fmt.Fprintf(os.Stderr, "%s: not enough arguments\n", progname)
    usage();
    os.Exit(1);
  }

  /* handle 1st arg - start page */
  s1 = os.Args[1]    /* !!! PBO */
  if len(s1) < 2 || s1[:2] != "-s" {
    fmt.Fprintf(os.Stderr, "%s: 1st arg should be -sstart_page\n", progname)
    usage();
    os.Exit(2);
  }
  i, _ = strconv.Atoi(s1[2:])
  if i < 1 || i > 65536 {
    fmt.Fprintf(os.Stderr, "%s: invalid start page %s\n", progname, s1[2:])
    usage();
    os.Exit(3);
  }
  psa.start_page = i

  /* handle 2nd arg - start page */
  s1 = os.Args[2]    /* !!! PBO */
  if len(s1) < 0 || s1[:2] != "-e" {
    fmt.Fprintf(os.Stderr, "%s: 2nd arg should be -eend_page\n", progname)
    usage();
    os.Exit(4);
  }
  i, _ = strconv.Atoi(s1[2:])
  if i < 1 || i > 65536 || i < psa.start_page {
    fmt.Fprintf(os.Stderr, "%s: invalid end page %s\n", progname, s1[2:])
    os.Exit(5);
  }
  psa.end_page = i

  argno = 3
  for _, s1 = range os.Args[3:] {
    if s1[0] != '-' { break }
    argno++
    switch s1[1] {
    case 'l':
      i, _ = strconv.Atoi(s1[2:])
      if i < 1 || i > 65536 {
        fmt.Fprintf(os.Stderr, "%s: invalid page length %s\n",
        progname, s1[2:])
        usage()
        os.Exit(6);
      }
      psa.page_len = i
    case 'f':
      if s1 != "-f" {
        fmt.Fprintf(os.Stderr, "%s: option should be \"-f\"\n",
        progname)
        usage()
        os.Exit(7);
      }
      psa.page_type = 'f'
    case 'd':
      if s1 == "-d" {
        fmt.Fprintf(os.Stderr, "%s: -d option requires a printer destination\n",
        progname)
        usage()
        os.Exit(8);
      }
      psa.print_dest = s1[2:]
    default:
      fmt.Fprintf(os.Stderr, "%s: unknown option %s\n",
      progname, s1)
      usage()
      os.Exit(9);
    }
  }

  if argno < len(os.Args) {
    s1 = os.Args[argno]
    psa.in_filename = s1
    if _, err := os.Stat(s1); err != nil && os.IsNotExist(err) {
      fmt.Fprintf(os.Stderr, "%s: input file \"%s\" does not exist\n",
      progname, s1)
      os.Exit(10);
    }
  } 
}


/*================================= process_input() ===============*/

const LINE_SIZE int = 10

func process_input(sa sp_args) {
  fin := os.Stdin /* input stream */
  fout := os.Stdout /* output stream */
  var s1 string
  var crc string
  var c int
  var line  string
  var line_ctr int
  var page_ctr int
  var inbuf string
  var err error

// in order to solve "declared and not used" problem
  s1 = s1
  crc = crc
  c = c
  line = line
  inbuf = inbuf
  err = err 


  if len(sa.in_filename) == 0 {
    fin = os.Stdin
  } else {
    fin, err = os.OpenFile(sa.in_filename, os.O_RDONLY, 0)
    if os.IsNotExist(err) {
      fmt.Fprintf(os.Stderr, "%s: could not open input file \"%s\"\n", progname, sa.in_filename)
      os.Exit(12);
    }
  }

  // ??? 
  if len(sa.print_dest) == 0 {
    fout = os.Stdout
  } else {
    fflush := bufio.NewWriter(os.Stdout) 
    fflush.Flush() 
    cmd := exec.Command("lp", "-d"+sa.print_dest) 
    _, err := cmd.Output() 
    if (err != nil) { 
      fmt.Fprintln(os.Stderr, progname, ": could not open pipe to \"lp -d", sa.print_dest,"\"") 
      os.Exit(13) 
    } 
  }

  /* begin one of two main loops based on page type */
  if sa.page_type == 'l' {
    line_ctr = 0;
    page_ctr = 1;

    for true {
      crc := make([]byte, LINE_SIZE)
      _, err := bufio.NewReaderSize(fin, INBUFSIZ).Read(crc)

      if err == nil {  /* error or EOF */
        break;
      }
      line_ctr++
      if line_ctr > sa.page_len {
        page_ctr++
        line_ctr = 1
      }
      if (page_ctr >= sa.start_page) && (page_ctr <= sa.end_page) {
        bufio.NewWriterSize(fout, INBUFSIZ).WriteString(string(crc))
      }
    }
  } else {
    page_ctr = 1
    for true {
      c, _, err := bufio.NewReaderSize(fin, INBUFSIZ).ReadRune()
      if err == io.EOF { /* error or EOF */
          break
      }
      if c == '\f' {  /* form feed */
        page_ctr++
      }
      if (page_ctr >= sa.start_page) && (page_ctr <= sa.end_page) {
        // output
        bufio.NewWriterSize(fout, INBUFSIZ).WriteRune(c)
      }
    }
  }

  /* end main loop */
  if page_ctr < sa.start_page {
    fmt.Fprintf(os.Stderr,
      "%s: start_page (%d) greater than total pages (%d)," +
      " no output written\n", progname, sa.start_page, page_ctr)
  } else if page_ctr < sa.end_page {
    fmt.Fprintf(os.Stderr,
      "%s: end_page (%d) greater than total pages (%d)," +
      " less output than expected\n", progname, sa.end_page, page_ctr)
  } else {
    fin.Close() 
    bufio.NewWriter(fout).Flush() 
    if len(sa.print_dest) != 0 { 
      fout.Close() 
    } 
    fmt.Fprintln(os.Stderr, progname, ": done") 
  }
}

 /* else if {
    fmt.Fprintf(os.Stderr,
      "%s: system error [%s] occurred on input stream fin\n", progname, s1);
    fin.Close()
    os.Exit(14) 
  }*/ 


/*================================= usage() =======================*/

func usage() {
  fmt.Fprintln(os.Stderr, "\nUSAGE: ", progname, " -sstart_page -eend_page [ -f | -llines_per_page ] [ -ddest ] [ in_filename ]")
}

/*================================= EOF ===========================*/

