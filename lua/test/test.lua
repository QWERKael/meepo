---
--- Generated by EmmyLua(https://github.com/EmmyLua)
--- Created by alex.
--- DateTime: 2021/5/12 下午3:52
---

local str = [[2021-05-12T14:57:56.961112Z 0 [Warning] No existing UUID has been found, so we assume that this is the first time that this server has been started. Generating a new UUID: 6fc1569a-b332-11eb-9b17-000c29809b1f.
2021-05-12T14:57:56.962661Z 0 [Warning] Gtid table is not ready to be used. Table 'mysql.gtid_executed' cannot be opened.
2021-05-12T14:57:57.390762Z 0 [Warning] CA certificate ca.pem is self signed.
2021-05-12T14:57:57.582669Z 1 [Note] A temporary password is generated for root@localhost: qYGen)B1&>(T
2021-05-12T14:58:01.377655Z 0 [Warning] TIMESTAMP with implicit DEFAULT value is deprecated. Please use --explicit_defaults_for_timestamp server option (see documentation for more details).
2021-05-12T14:58:01.379445Z 0 [Note] /usr/sbin/mysqld (mysqld 5.7.33-36) starting as process 55768 ...
2021-05-12T14:58:01.383337Z 0 [Note] InnoDB: PUNCH HOLE support available]]
--for word in string.gmatch(s, "[%S]+") do
--for word in string.gmatch(s, "\n[%S]+") do
--    print(word)
--end

local pattern = "A temporary password is generated for root@localhost: ([%S]+)\n"
local _, _, cap1 = string.find(str, pattern)
print(cap1)

--local stringx = require("pl.stringx")
--local pretty = require("pl.pretty")
--local lines = stringx.splitlines(s)
--pretty.dump(lines)
--local names = {}
--for idx, line in ipairs(lines) do
--    local n = stringx.split(line)
--    print(n[1])
--    names[idx] = n[1]
--end
--
--pretty.dump(names)
--
--local args = stringx.join((' '), names)
--print("yum remove -y "..args)