// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build tools
// +build tools

package sole

// example to get latest .git-subtree.sh
//go:generate bash -c "curl -s -L -o .git-subtree.sh https://gist.githubusercontent.com/searKing/8e948af12f03074b1a7c07e1cba2c407/raw/f3d581915f5338e058b20084c5c8850fcfcaf462/.git-subtree.sh"
//go:generate bash -c "chmod a+x .git-subtree.sh"
//go:generate bash -c "./.git-subtree.sh purge --remove_files"
//go:generate bash -c "./.git-subtree.sh tidy -u"
