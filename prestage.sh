# Copyright 2016 Comcast Cable Communications Management, LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

#! /bin/bash

# because the seemingly only sane way to call this script is with parallel, we have to sleep for a bit to let Zookeeper
# come up properly.
sleep 10

echo "prestaging data."

zkCli.sh localhost:2181 -cmd create /$BUCKET_1
zkCli.sh localhost:2181 -cmd create /$BUCKET_2

exit