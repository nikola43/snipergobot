// SPDX-License-Identifier: MIT
pragma solidity ^0.8.13;

contract Buy {
    
    event Received(address, uint, bytes);
    event FallbackEvent(address, uint, bytes);

    receive() external payable {
        processPayment();
    }

    fallback() external payable {
        processFallback();
    }

    function processFallback() private {
      emit FallbackEvent(msg.sender, msg.value, msg.data);
    }

    function processPayment() private {
      emit Received(msg.sender, msg.value, msg.data);
    }
}

