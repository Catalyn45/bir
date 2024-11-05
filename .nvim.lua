require('dap-go').setup {
    delve = {
        detached = false,
    },
	dap_configurations = {
		{
			type = 'go',
			name = 'Debug package args',
			request = 'launch',
			program = '${fileDirname}',
			args = {"./test.bir"},
		},
	},
}
