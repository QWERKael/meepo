---
--- Generated by EmmyLua(https://github.com/EmmyLua)
--- Created by chenyixuan.
--- DateTime: 2021/5/7 1:36 下午
---

function main()
    db = require('db')
    local mysql, err = db.open("mysql", "root:1234@(localhost:3307)/pitaya_game?charset=utf8&parseTime=True&loc=Local")
    if err then
        error(err)
    end

end