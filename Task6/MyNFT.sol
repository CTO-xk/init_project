// SPDX-License-Identifier: MIT
pragma solidity ^0.8;

// 导入OpenZeppelin的核心库
import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/token/ERC721/extensions/ERC721URIStorage.sol";
import "@openzeppelin/contracts/utils/Counters.sol";

/**
 * @title MyNFT
 * @dev 基于ERC721标准的NFT合约，支持铸造和元数据管理
 */
contract MyNFT is ERC721URIStorage {
    // 使用计数器自动生成唯一tokenID
    using Counters for Counters.Counter;
    Counters.Counter private _tokenIdCounter;

    // 合约所有者地址
    address private immutable _owner;

    // 仅所有者可调用的修饰符
    modifier onlyOwner() {
        require(msg.sender == _owner, "MyNFT: caller is not the owner");
        _;
    }

    /**
     * @dev 构造函数，初始化NFT名称和符号
     * @param name NFT集合名称
     * @param symbol NFT符号
     */
    constructor(string memory name, string memory symbol) ERC721(name, symbol) {
        _owner = msg.sender; // 部署者成为所有者
    }

    /**
     * @dev 铸造新NFT并关联元数据
     * @param recipient 接收NFT的地址
     * @param _tokenURI IPFS上的元数据链接
     * @return 新铸造NFT的tokenID
     */
    function mintNFT(address recipient, string memory _tokenURI) 
        public 
        onlyOwner 
        returns (uint256) 
    {
        // 验证接收地址有效
        require(recipient != address(0), "MyNFT: recipient is zero address");
        // 验证元数据链接不为空
        require(bytes(_tokenURI).length > 0, "MyNFT: tokenURI is empty");

        // 获取当前tokenID并自增
        uint256 tokenId = _tokenIdCounter.current();
        _tokenIdCounter.increment();

        // 安全铸造NFT并分配给接收者
        _safeMint(recipient, tokenId);
        // 关联元数据链接
        _setTokenURI(tokenId, _tokenURI);

        return tokenId;
    }

    /**
     * @dev 重写tokenURI函数，确保兼容性
     */
    function tokenURI(uint256 tokenId)
        public
        view
        override(ERC721URIStorage)
        returns (string memory)
    {
        return super.tokenURI(tokenId);
    }

    /**
     * @dev 获取合约所有者地址
     */
    function owner() public view returns (address) {
        return _owner;
    }
}
