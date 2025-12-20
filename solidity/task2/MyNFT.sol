// SPDX-License-Identifier: MIT
pragma solidity ^0.8.21;
import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/utils/Strings.sol";

contract MyNFT is ERC721 {
    uint256 counter;
    using Strings for uint256;
    mapping(uint256 => string) private _tokenURIs;

    constructor(
        string memory name_,
        string memory symbol_
    ) ERC721(name_, symbol_) {}

    function mintNFT(
        address to,
        string memory _tokenURI
    ) public returns (uint256) {
        super._safeMint(to, counter++);
        _tokenURIs[counter] = _tokenURI;
        return counter;
    }

    function tokenURI(
        uint256 tokenId
    ) public view override(ERC721) returns (string memory) {
        _requireOwned(tokenId);

        string memory baseURI = _baseURI();
        return
            bytes(baseURI).length > 0
                ? string.concat(baseURI, _tokenURIs[counter])
                : "";
    }

    function _baseURI() internal pure override(ERC721) returns (string memory) {
        // IPFS 地址
        return "https://maroon-decent-dragonfly-611.mypinata.cloud/ipfs/";
    }
}
