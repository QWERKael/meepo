---
--- Generated by EmmyLua(https://github.com/EmmyLua)
--- Created by chenyixuan.
--- DateTime: 2021/4/22 10:46 上午
---

local task_string = [[{
 "PluginName": "show",
 "Command": "info",
 "Args": {
  "SubCommands": null,
  "Flags": null,
  "KVs": null
 },
 "ExtraArgs": null,
 "Result": null
}]]

--print("Lua任务内容：", task_string)

print("+++++++++++++++++++++++++++++++++++++++++++++++++++++")
local uuid = send("127.0.0.1:4002", "server", task_string)
print("UUID is : ", uuid)
print("Result is : ", getRst(uuid))
print("+++++++++++++++++++++++++++++++++++++++++++++++++++++")
local uuid = send("127.0.0.1:4003", "server", task_string)
print("UUID is : ", uuid)
print("Result is : ", getRst(uuid))
print("+++++++++++++++++++++++++++++++++++++++++++++++++++++")
local uuid = send("127.0.0.1:4004", "server", task_string)
print("UUID is : ", uuid)
print("Result is : ", getRst(uuid))
print("+++++++++++++++++++++++++++++++++++++++++++++++++++++")