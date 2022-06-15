
LINUX=false
OS="$(uname)"
if [[ "${OS}" == "Linux" ]]; then
  LINUX=true
fi

UNAME_MACHINE="$(/usr/bin/uname -m)"

if [[ UNAME_MACHINE == "x86_64" && LINUX ]]; then
