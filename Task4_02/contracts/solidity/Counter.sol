// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

contract Counter {
    uint256 private _count;
    
    // 当计数更新时触发的事件
    event CountUpdated(uint256 newCount);
    
    // 构造函数，初始化计数
    constructor(uint256 initialCount) {
        _count = initialCount;
        emit CountUpdated(_count);
    }
    
    // 获取当前计数
    function getCount() public view returns (uint256) {
        return _count;
    }
    
    // 增加计数
    function increment() public {
        _count += 1;
        emit CountUpdated(_count);
    }
    
    // 减少计数
    function decrement() public {
        require(_count > 0, "Count cannot be negative");
        _count -= 1;
        emit CountUpdated(_count);
    }
    
    // 重置计数为0
    function reset() public {
        _count = 0;
        emit CountUpdated(_count);
    }
}
