// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract UnsafeBank {
    mapping(address => uint256) public balances;
    bool private locked;

    // Event to log successful withdrawals for better traceability
    event Withdrawn(address indexed user, uint256 amount);

    modifier noReentrancy() {
        require(!locked, "Reentrant call");
        locked = true;
        _;
        locked = false;
    }

    function deposit() external payable {
        balances[msg.sender] += msg.value;
    }

    function withdraw() external noReentrancy {
        // Cache the balance to avoid multiple storage reads
        uint256 amount = balances[msg.sender];

        // Validate sufficient balance (improved error message for clarity)
        require(amount > 0, "Insufficient balance: no funds available");

        // Apply Checks-Effects-Interactions pattern: update state before external interaction
        balances[msg.sender] = 0;

        // Perform the external ether transfer using low-level call for flexibility
        // Cast msg.sender to payable for explicit type safety
        (bool success, ) = payable(msg.sender).call{value: amount}("");

        // Revert on transfer failure to maintain atomicity
        require(success, "Ether transfer failed: unable to send funds");

        // Emit event for off-chain monitoring and debugging
        emit Withdrawn(msg.sender, amount);
    }
}
