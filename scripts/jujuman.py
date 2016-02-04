"""Functions for generating the manpage using the juju command."""

import subprocess
import textwrap
import time


class JujuMan(object):

    def __init__(self):
        self.version = self._version()

    def get_filename(self, options):
        """Provides name of manpage"""
        return 'juju.1'

    def run_juju(self, *args):
        cmd = ['juju'] + list(args)
        return subprocess.check_output(cmd).strip()

    def _version(self):
        juju_version = self.run_juju('version')
        return juju_version.split('-')[0]

    def commands(self):
        commands = self.run_juju('help', 'commands')
        result = []
        for line in commands.split('\n'):
            name, short_help = line.split(' ', 1)
            if 'alias for' in short_help:
                continue
            result.append((name, short_help.strip()))
        return result

    def write_documentation(self, options, outfile):
        """Assembles a man page"""
        t = time.time()
        tt = time.gmtime(t)
        params = {
            "cmd": "juju",
            "datestamp": time.strftime("%Y-%m-%d",tt),
            "timestamp": time.strftime("%Y-%m-%d %H:%M:%S +0000",tt),
            "version": self.version,
            }
        outfile.write(man_preamble % params)
        outfile.write(man_escape(man_head % params))
        outfile.write(man_escape(self.getcommand_list(params)))
        outfile.write("".join(environment_variables()))
        outfile.write(man_escape(man_foot % params))

    def getcommand_list(self, params):
        """Builds summary help for command names in manpage format"""
        output = '.SH "COMMAND OVERVIEW"\n'
        for cmd_name, short_help in self.commands():
            tmp = '.TP\n.B "%s %s"\n%s\n' % (params['cmd'], cmd_name, short_help)
            output = output + tmp
        return output


ENVIRONMENT = (
    ('JUJU_MODEL', textwrap.dedent("""\
        Provides a way for the shell environment to specify the current Juju
        model to use.  If the model is specified explicitly using
        -m MODEL, this takes precedence.
        """)),
    ('JUJU_HOME', textwrap.dedent("""\
        Overrides the default Juju configuration directory of ~/.juju.
        """)),
    ('AWS_ACCESS_KEY_ID', textwrap.dedent("""\
        The access-key for your AWS account.
        """)),
    ('AWS_SECRET_ACCESS_KEY', textwrap.dedent("""\
        The secret-key for your AWS account.
        """)),
    ('OS_USERNAME', textwrap.dedent("""\
        Your openstack username.
        """)),
    ('OS_PASSWORD', textwrap.dedent("""\
        Your openstack password.
        """)),
    ('OS_TENANT_NAME', textwrap.dedent("""\
        Your openstack tenant name.
        """)),
    ('OS_REGION_NAME', textwrap.dedent("""\
        Your openstack region name.
        """)),
)

def man_escape(string):
    """Escapes strings for man page compatibility"""
    result = string.replace("\\","\\\\")
    result = result.replace("`","\\'")
    result = result.replace("'","\\*(Aq")
    result = result.replace("-","\\-")
    return result


def environment_variables():
    yield ".SH \"ENVIRONMENT\"\n"

    for k, desc in ENVIRONMENT:
        yield ".TP\n"
        yield ".I \"%s\"\n" % k
        yield man_escape(desc) + "\n"


man_preamble = """\
.\\\"Man page for Juju (%(cmd)s)
.\\\"
.\\\" Large parts of this file are autogenerated from the output of
.\\\"     \"%(cmd)s help commands\"
.\\\"     \"%(cmd)s help <cmd>\"
.\\\"
.\\\" Generation time: %(timestamp)s
.\\\"

.ie \\n(.g .ds Aq \\(aq
.el .ds Aq '
"""

man_head = """\
.TH %(cmd)s 1 "%(datestamp)s" "%(version)s" "Juju"
.SH "NAME"
%(cmd)s - Juju -- devops distilled
.SH "SYNOPSIS"
.B "%(cmd)s"
.I "command"
[
.I "command_options"
]
.br
.B "%(cmd)s"
.B "help"
.br
.B "%(cmd)s"
.B "help"
.I "command"
.SH "DESCRIPTION"

Juju provides easy, intelligent service orchestration on top of environments
such as OpenStack, Amazon AWS, or bare metal.
"""

man_foot = """\
.SH "FILES"
.TP
.I "~/.juju/environments.yaml"
This is the Juju config file, which you can use to specify multiple
models in which to deploy.

A config file can be created using
.B juju init
which you can then edit to provide the secret keys, or use environment
variables to provide the secret values.

.SH "SEE ALSO"
.UR https://juju.ubuntu.com/
.BR https://juju.ubuntu.com/
"""

