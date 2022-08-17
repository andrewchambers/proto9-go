package proto9

import (
	"fmt"
)

func (e *Rlerror) Error() string {
	if 0 <= int(e.Ecode) && int(e.Ecode) < len(dotlErrors) {
		s := dotlErrors[e.Ecode]
		if s != "" {
			return s
		}
	}
	return fmt.Sprintf("Error: errno(%d)", e.Ecode)
}

// numbers defined on Linux/amd64.
const (
	E2BIG           = 0x7
	EACCES          = 0xd
	EADDRINUSE      = 0x62
	EADDRNOTAVAIL   = 0x63
	EADV            = 0x44
	EAFNOSUPPORT    = 0x61
	EAGAIN          = 0xb
	EALREADY        = 0x72
	EBADE           = 0x34
	EBADF           = 0x9
	EBADFD          = 0x4d
	EBADMSG         = 0x4a
	EBADR           = 0x35
	EBADRQC         = 0x38
	EBADSLT         = 0x39
	EBFONT          = 0x3b
	EBUSY           = 0x10
	ECANCELED       = 0x7d
	ECHILD          = 0xa
	ECHRNG          = 0x2c
	ECOMM           = 0x46
	ECONNABORTED    = 0x67
	ECONNREFUSED    = 0x6f
	ECONNRESET      = 0x68
	EDEADLK         = 0x23
	EDEADLOCK       = 0x23
	EDESTADDRREQ    = 0x59
	EDOM            = 0x21
	EDOTDOT         = 0x49
	EDQUOT          = 0x7a
	EEXIST          = 0x11
	EFAULT          = 0xe
	EFBIG           = 0x1b
	EHOSTDOWN       = 0x70
	EHOSTUNREACH    = 0x71
	EHWPOISON       = 0x85
	EIDRM           = 0x2b
	EILSEQ          = 0x54
	EINPROGRESS     = 0x73
	EINTR           = 0x4
	EINVAL          = 0x16
	EIO             = 0x5
	EISCONN         = 0x6a
	EISDIR          = 0x15
	EISNAM          = 0x78
	EKEYEXPIRED     = 0x7f
	EKEYREJECTED    = 0x81
	EKEYREVOKED     = 0x80
	EL2HLT          = 0x33
	EL2NSYNC        = 0x2d
	EL3HLT          = 0x2e
	EL3RST          = 0x2f
	ELIBACC         = 0x4f
	ELIBBAD         = 0x50
	ELIBEXEC        = 0x53
	ELIBMAX         = 0x52
	ELIBSCN         = 0x51
	ELNRNG          = 0x30
	ELOOP           = 0x28
	EMEDIUMTYPE     = 0x7c
	EMFILE          = 0x18
	EMLINK          = 0x1f
	EMSGSIZE        = 0x5a
	EMULTIHOP       = 0x48
	ENAMETOOLONG    = 0x24
	ENAVAIL         = 0x77
	ENETDOWN        = 0x64
	ENETRESET       = 0x66
	ENETUNREACH     = 0x65
	ENFILE          = 0x17
	ENOANO          = 0x37
	ENOBUFS         = 0x69
	ENOCSI          = 0x32
	ENODATA         = 0x3d
	ENODEV          = 0x13
	ENOENT          = 0x2
	ENOEXEC         = 0x8
	ENOKEY          = 0x7e
	ENOLCK          = 0x25
	ENOLINK         = 0x43
	ENOMEDIUM       = 0x7b
	ENOMEM          = 0xc
	ENOMSG          = 0x2a
	ENONET          = 0x40
	ENOPKG          = 0x41
	ENOPROTOOPT     = 0x5c
	ENOSPC          = 0x1c
	ENOSR           = 0x3f
	ENOSTR          = 0x3c
	ENOSYS          = 0x26
	ENOTBLK         = 0xf
	ENOTCONN        = 0x6b
	ENOTDIR         = 0x14
	ENOTEMPTY       = 0x27
	ENOTNAM         = 0x76
	ENOTRECOVERABLE = 0x83
	ENOTSOCK        = 0x58
	ENOTSUP         = 0x5f
	ENOTTY          = 0x19
	ENOTUNIQ        = 0x4c
	ENXIO           = 0x6
	EOPNOTSUPP      = 0x5f
	EOVERFLOW       = 0x4b
	EOWNERDEAD      = 0x82
	EPERM           = 0x1
	EPFNOSUPPORT    = 0x60
	EPIPE           = 0x20
	EPROTO          = 0x47
	EPROTONOSUPPORT = 0x5d
	EPROTOTYPE      = 0x5b
	ERANGE          = 0x22
	EREMCHG         = 0x4e
	EREMOTE         = 0x42
	EREMOTEIO       = 0x79
	ERESTART        = 0x55
	ERFKILL         = 0x84
	EROFS           = 0x1e
	ESHUTDOWN       = 0x6c
	ESOCKTNOSUPPORT = 0x5e
	ESPIPE          = 0x1d
	ESRCH           = 0x3
	ESRMNT          = 0x45
	ESTALE          = 0x74
	ESTRPIPE        = 0x56
	ETIME           = 0x3e
	ETIMEDOUT       = 0x6e
	ETOOMANYREFS    = 0x6d
	ETXTBSY         = 0x1a
	EUCLEAN         = 0x75
	EUNATCH         = 0x31
	EUSERS          = 0x57
	EWOULDBLOCK     = 0xb
	EXDEV           = 0x12
	EXFULL          = 0x36
)

var dotlErrors = [...]string{
	1:   "operation not permitted",
	2:   "no such file or directory",
	3:   "no such process",
	4:   "interrupted system call",
	5:   "input/output error",
	6:   "no such device or address",
	7:   "argument list too long",
	8:   "exec format error",
	9:   "bad file descriptor",
	10:  "no child processes",
	11:  "resource temporarily unavailable",
	12:  "cannot allocate memory",
	13:  "permission denied",
	14:  "bad address",
	15:  "block device required",
	16:  "device or resource busy",
	17:  "file exists",
	18:  "invalid cross-device link",
	19:  "no such device",
	20:  "not a directory",
	21:  "is a directory",
	22:  "invalid argument",
	23:  "too many open files in system",
	24:  "too many open files",
	25:  "inappropriate ioctl for device",
	26:  "text file busy",
	27:  "file too large",
	28:  "no space left on device",
	29:  "illegal seek",
	30:  "read-only file system",
	31:  "too many links",
	32:  "broken pipe",
	33:  "numerical argument out of domain",
	34:  "numerical result out of range",
	35:  "resource deadlock avoided",
	36:  "file name too long",
	37:  "no locks available",
	38:  "function not implemented",
	39:  "directory not empty",
	40:  "too many levels of symbolic links",
	42:  "no message of desired type",
	43:  "identifier removed",
	44:  "channel number out of range",
	45:  "level 2 not synchronized",
	46:  "level 3 halted",
	47:  "level 3 reset",
	48:  "link number out of range",
	49:  "protocol driver not attached",
	50:  "no CSI structure available",
	51:  "level 2 halted",
	52:  "invalid exchange",
	53:  "invalid request descriptor",
	54:  "exchange full",
	55:  "no anode",
	56:  "invalid request code",
	57:  "invalid slot",
	59:  "bad font file format",
	60:  "device not a stream",
	61:  "no data available",
	62:  "timer expired",
	63:  "out of streams resources",
	64:  "machine is not on the network",
	65:  "package not installed",
	66:  "object is remote",
	67:  "link has been severed",
	68:  "advertise error",
	69:  "srmount error",
	70:  "communication error on send",
	71:  "protocol error",
	72:  "multihop attempted",
	73:  "RFS specific error",
	74:  "bad message",
	75:  "value too large for defined data type",
	76:  "name not unique on network",
	77:  "file descriptor in bad state",
	78:  "remote address changed",
	79:  "can not access a needed shared library",
	80:  "accessing a corrupted shared library",
	81:  ".lib section in a.out corrupted",
	82:  "attempting to link in too many shared libraries",
	83:  "cannot exec a shared library directly",
	84:  "invalid or incomplete multibyte or wide character",
	85:  "interrupted system call should be restarted",
	86:  "streams pipe error",
	87:  "too many users",
	88:  "socket operation on non-socket",
	89:  "destination address required",
	90:  "message too long",
	91:  "protocol wrong type for socket",
	92:  "protocol not available",
	93:  "protocol not supported",
	94:  "socket type not supported",
	95:  "operation not supported",
	96:  "protocol family not supported",
	97:  "address family not supported by protocol",
	98:  "address already in use",
	99:  "cannot assign requested address",
	100: "network is down",
	101: "network is unreachable",
	102: "network dropped connection on reset",
	103: "software caused connection abort",
	104: "connection reset by peer",
	105: "no buffer space available",
	106: "transport endpoint is already connected",
	107: "transport endpoint is not connected",
	108: "cannot send after transport endpoint shutdown",
	109: "too many references: cannot splice",
	110: "connection timed out",
	111: "connection refused",
	112: "host is down",
	113: "no route to host",
	114: "operation already in progress",
	115: "operation now in progress",
	116: "stale NFS file handle",
	117: "structure needs cleaning",
	118: "not a XENIX named type file",
	119: "no XENIX semaphores available",
	120: "is a named type file",
	121: "remote I/O error",
	122: "disk quota exceeded",
	123: "no medium found",
	124: "wrong medium type",
	125: "operation canceled",
	126: "required key not available",
	127: "key has expired",
	128: "key has been revoked",
	129: "key was rejected by service",
	130: "owner died",
	131: "state not recoverable",
	132: "operation not possible due to RF-kill",
}
