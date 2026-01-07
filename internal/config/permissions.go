package config

import "os"

// DirPerm | Directory permissions:
// owner: read/write/execute
// group: read/execute
// other: read/execute
const DirPerm = os.FileMode(0o755)

// FilePerm | File permissions:
// owner: read/write
// group: read
// other: read
const FilePerm = os.FileMode(0o644)
