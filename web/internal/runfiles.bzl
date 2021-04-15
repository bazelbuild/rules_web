# Copyright 2018 Google Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
"""Runfiles modules contains utility functions for working with runfiles."""

def _collect(ctx, files = [], targets = []):
    """Builds a runfiles object with transitive files from all targets.

    Args:
        ctx: Context object for the rule where this being used.
        files: a list of File object to include in the runfiles.
        targets: a list of Target object from which runfiles will be collected.

    Returns:
        A configured runfiles object that include data and default runfiles for the
        rule, all transitive runfiles from targets, and all files from files.
    """
    transitive_runfiles_list = []
    default_runfiles = []
    data_runfiles = []
    for target in targets:
        if hasattr(target, "transitive_runfiles"):
            transitive_runfiles_list.extend(target.transitive_runfiles)
        if hasattr(target, "default_runfiles"):
            default_runfiles.append(target.default_runfiles)
        if hasattr(target, "data_runfiles"):
            data_runfiles.append(target.data_runfiles)
        if hasattr(target, "files"):
            transitive_runfiles_list.append(target.files)

    transitive_runfiles = depset(transitive = transitive_runfiles_list)

    result = ctx.runfiles(
        collect_data = True,
        collect_default = True,
        files = files,
        transitive_files = transitive_runfiles,
    )

    for default in default_runfiles:
        result = result.merge(default)

    for data in data_runfiles:
        result = result.merge(data)

    return result

runfiles = struct(collect = _collect)
