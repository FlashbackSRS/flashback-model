#!/bin/bash
set -e

go test github.com/FlashbackSRS/flashback-model github.com/FlashbackSRS/flashback-model/test
gopherjs test --tags=disableunsafe github.com/FlashbackSRS/flashback-model github.com/FlashbackSRS/flashback-model/test
