#!/bin/sh
# Run the versioned CLI from a private cache, downloading its release if needed.
set -eu
umask 077
fail() {
    printf '%s\n' '{"schemaVersion":"v1","error_code":"BootstrapError","message":"Could not prepare the prebuilt Inspire CLI.","recovery":{"action":"inspect","hint":"Check GitHub access, archive tools and the matching release, or install inspire from its release archive."}}' >&2
    exit 1
}
script_dir=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd) || fail
version=$(tr -d '\r\n' < "$script_dir/VERSION")
case "$version" in ''|*[!0-9.]*|.*|*.) fail ;; esac
case "$(uname -s)" in Linux) target_os=linux ;; Darwin) target_os=darwin ;; *) fail ;; esac
case "$(uname -m)" in x86_64|amd64) target_arch=amd64 ;; arm64|aarch64) target_arch=arm64 ;; *) fail ;; esac
cache_root=${INSPIRE_CLI_CACHE_DIR:-${XDG_CACHE_HOME:-${HOME:?}/.cache}/inspire/cli}
cache_dir="$cache_root/$version/${target_os}_$target_arch"
binary="$cache_dir/inspire"
if [ ! -x "$binary" ]; then
    mkdir -p "$cache_dir" || fail
    staging=$(mktemp -d "$cache_dir/.download.XXXXXXXX") || fail
    trap 'rm -rf -- "$staging"' EXIT HUP INT TERM
    asset="inspire_${version}_${target_os}_${target_arch}.tar.gz"
    release="https://github.com/Lixsp11/inspire-go-sdk/releases/download/v$version"
    curl --fail --location --silent --show-error --proto '=https' --tlsv1.2 "$release/SHA256SUMS" -o "$staging/SHA256SUMS" 2>"$staging/download.log" || fail
    curl --fail --location --silent --show-error --proto '=https' --tlsv1.2 "$release/$asset" -o "$staging/$asset" 2>"$staging/download.log" || fail
    expected=$(awk -v name="$asset" '$2 == name {print $1}' "$staging/SHA256SUMS")
    [ "${#expected}" -eq 64 ] || fail
    case "$expected" in *[!0-9a-f]*) fail ;; esac
    if command -v sha256sum >/dev/null 2>&1; then
        actual=$(sha256sum "$staging/$asset" | awk '{print $1}')
    elif command -v shasum >/dev/null 2>&1; then
        actual=$(shasum -a 256 "$staging/$asset" | awk '{print $1}')
    else
        fail
    fi
    [ "$actual" = "$expected" ] || fail
    mkdir "$staging/payload" || fail
    tar -xzf "$staging/$asset" -C "$staging/payload" 2>"$staging/extract.log" || fail
    [ -f "$staging/payload/inspire" ] || fail
    cp "$staging/payload/LICENSE" "$staging/payload/NOTICE" "$cache_dir/" || fail
    cp -R "$staging/payload/third_party" "$cache_dir/" || fail
    chmod 700 "$staging/payload/inspire" || fail
    mv -f "$staging/payload/inspire" "$binary" || fail
    rm -rf -- "$staging"
    trap - EXIT HUP INT TERM
fi
signature=$(od -An -tx1 -N4 "$binary" 2>/dev/null | tr -d ' \n') || fail
case "$target_os:$signature" in linux:7f454c46|darwin:cffaedfe) ;; *) fail ;; esac
error_file=$(mktemp "$cache_dir/.stderr.XXXXXXXX") || fail
trap 'rm -f -- "$error_file"' EXIT
if (exec "$binary" "$@") 2>"$error_file"; then
    cli_status=0
else
    cli_status=$?
fi
if [ "$cli_status" -eq 126 ] || [ "$cli_status" -eq 127 ]; then fail; fi
cat "$error_file" >&2
exit "$cli_status"
