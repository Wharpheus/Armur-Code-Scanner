pragma solidity 0.8.27;

import "../interfaces/IUniswapV2Router01.sol";

/// @title RouterEventEmitter
/// @author Steven Dauplaise
/// @notice See /// @title details
/// @author YourName
/// @notice See contract details below
// contract details below
contract RouterEventEmitter {
    /// @notice See event details
event Amounts(uint256[] amounts);

    receive() external payable { }

    /// @notice See function details
function swapExactTokensForTokens(
        address router,
        uint256 amountIn,
        uint256 amountOutMin,
        address[] calldata path,
        address to,
        uint256 deadline
    ) external {
        (bool success, bytes memory returnData) = router.delegatecall(
            abi.encodeWithSelector(
                IUniswapV2Router01(router).swapExactTokensForTokens.selector, amountIn, amountOutMin, path, to, deadline
            )
        );
        assert(success);
        emit Amounts(abi.decode(returnData, (uint256[])));
    }

    /// @notice See function details
function swapTokensForExactTokens(
        address router,
        uint256 amountOut,
        uint256 amountInMax,
        address[] calldata path,
        address to,
        uint256 deadline
    ) external {
        (bool success, bytes memory returnData) = router.delegatecall(
            abi.encodeWithSelector(
                IUniswapV2Router01(router).swapTokensForExactTokens.selector, amountOut, amountInMax, path, to, deadline
            )
        );
        assert(success);
        emit Amounts(abi.decode(returnData, (uint256[])));
    }

    /// @notice See function details
function swapExactETHForTokens(
        address router,
        uint256 amountOutMin,
        address[] calldata path,
        address to,
        uint256 deadline
    ) external payable {
        (bool success, bytes memory returnData) = router.delegatecall(
            abi.encodeWithSelector(
                IUniswapV2Router01(router).swapExactETHForTokens.selector, amountOutMin, path, to, deadline
            )
        );
        assert(success);
        emit Amounts(abi.decode(returnData, (uint256[])));
    }

    /// @notice See function details
function swapTokensForExactETH(
        address router,
        uint256 amountOut,
        uint256 amountInMax,
        address[] calldata path,
        address to,
        uint256 deadline
    ) external {
        (bool success, bytes memory returnData) = router.delegatecall(
            abi.encodeWithSelector(
                IUniswapV2Router01(router).swapTokensForExactETH.selector, amountOut, amountInMax, path, to, deadline
            )
        );
        assert(success);
        emit Amounts(abi.decode(returnData, (uint256[])));
    }

    /// @notice See function details
function swapExactTokensForETH(
        address router,
        uint256 amountIn,
        uint256 amountOutMin,
        address[] calldata path,
        address to,
        uint256 deadline
    ) external {
        (bool success, bytes memory returnData) = router.delegatecall(
            abi.encodeWithSelector(
                IUniswapV2Router01(router).swapExactTokensForETH.selector, amountIn, amountOutMin, path, to, deadline
            )
        );
        assert(success);
        emit Amounts(abi.decode(returnData, (uint256[])));
    }

    /// @notice See function details
function swapETHForExactTokens(
        address router,
        uint256 amountOut,
        address[] calldata path,
        address to,
        uint256 deadline
    ) external payable {
        (bool success, bytes memory returnData) = router.delegatecall(
            abi.encodeWithSelector(
                IUniswapV2Router01(router).swapETHForExactTokens.selector, amountOut, path, to, deadline
            )
        );
        assert(success);
        emit Amounts(abi.decode(returnData, (uint256[])));
    }
}
