# Copyright 2016 The Closure Rules Authors. All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

_PASS_PROPS = (
    "neverlink",
    "testonly_",
    "visibility",
    "exports",
    "runtime_deps",
    "deps",
)


def _java_import_external(repository_ctx):
  """Downloads jar and creates a java_import rule."""
  if (repository_ctx.attr.generated_linkable_rule_name and
      not repository_ctx.attr.neverlink):
    fail("Only use generated_linkable_rule_name if neverlink is set")
  name = repository_ctx.attr.generated_rule_name or repository_ctx.name
  urls = repository_ctx.attr.jar_urls
  sha = repository_ctx.attr.jar_sha256
  path = repository_ctx.name + ".jar"
  for url in urls:
    if url.endswith(".jar"):
      path = url[url.rindex("/") + 1:]
      break
  srcurls = repository_ctx.attr.srcjar_urls
  srcsha = repository_ctx.attr.srcjar_sha256
  srcpath = repository_ctx.name + "-src.jar" if srcurls else ""
  for url in srcurls:
    if url.endswith(".jar"):
      srcpath = url[url.rindex("/") + 1:].replace("-sources.jar", "-src.jar")
      break
  lines = ["# DO NOT EDIT: generated by java_import_external()", ""]
  if repository_ctx.attr.default_visibility:
    lines.append("package(default_visibility = %s)" %
                 (repository_ctx.attr.default_visibility))
    lines.append("")
  lines.append("licenses(%s)" % repr(repository_ctx.attr.licenses))
  lines.append("")
  lines.extend(
      _make_java_import(name, path, srcpath, repository_ctx.attr, _PASS_PROPS))
  if (repository_ctx.attr.neverlink and
      repository_ctx.attr.generated_linkable_rule_name):
    lines.extend(
        _make_java_import(repository_ctx.attr.generated_linkable_rule_name,
                          path, srcpath, repository_ctx.attr,
                          [p for p in _PASS_PROPS if p != "neverlink"]))
  extra = repository_ctx.attr.extra_build_file_content
  if extra:
    lines.append(extra)
    if not extra.endswith("\n"):
      lines.append("")
  repository_ctx.download(urls, path, sha)
  if srcurls:
    repository_ctx.download(srcurls, srcpath, srcsha)
  repository_ctx.file("BUILD", "\n".join(lines))


def _make_java_import(name, path, srcpath, attrs, props):
  lines = [
      "java_import(",
      "    name = %s," % repr(name),
      "    jars = [%s]," % repr(path),
  ]
  if srcpath:
    lines.append("    srcjar = %s," % repr(srcpath))
  for prop in props:
    value = getattr(attrs, prop, None)
    if value:
      if prop.endswith("_"):
        prop = prop[:-1]
      lines.append("    %s = %s," % (prop, repr(value)))
  lines.append(")")
  lines.append("")
  return lines


java_import_external = repository_rule(
    implementation=_java_import_external,
    attrs={
        "licenses":
            attr.string_list(mandatory=True, allow_empty=False),
        "jar_urls":
            attr.string_list(mandatory=True, allow_empty=False),
        "jar_sha256":
            attr.string(mandatory=True),
        "srcjar_urls":
            attr.string_list(),
        "srcjar_sha256":
            attr.string(),
        "deps":
            attr.string_list(),
        "runtime_deps":
            attr.string_list(),
        "testonly_":
            attr.bool(),
        "exports":
            attr.string_list(),
        "neverlink":
            attr.bool(),
        "generated_rule_name":
            attr.string(),
        "generated_linkable_rule_name":
            attr.string(),
        "default_visibility":
            attr.string_list(default=["//visibility:public"]),
        "extra_build_file_content":
            attr.string(),
    })
