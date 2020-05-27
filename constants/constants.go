package constants

const(

	// https://gist.github.com/asukakenji/f15ba7e588ac42795f421b48b8aede63

	OS_LINUX   = "linux"
	OS_WINDOWS = "windows"
	OS_MAC     = "darwin"

	ARCH_32bit       = "386"
	ARCH_64bit       = "amd64"
	ARCH_ARM_32bit   = "arm"
	ARCH_ARM_64bit   = "arm64"


	REGEX_SEMANTIC_VER   = `(\d+\.)?(\d+\.)?(\*|\d+)`
	REGEX_CURRENT_VER    = `v`+REGEX_SEMANTIC_VER+`$`
	REGEX_NEXT_VER       = REGEX_SEMANTIC_VER+`[.]\s`

	TERRAFORM_DOWNLOAD_BASE_URL = "https://releases.hashicorp.com/terraform/"
)
