#!/bin/env bash

B64=$(<<<$1 base64 | tr -d '\n')
LDFLAGS=$(<<<$1 sed -e 's:%:%%:g')

echo ": *.go |> GOOS=@(OS) GOARCH=@(ARCH) go build -o %o -ldflags \"${LDFLAGS} -X main.BUILD_COMMAND_B64=${B64}\" %f |> @(BINARY)"
