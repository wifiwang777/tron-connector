// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

interface IERC20 {
    function transfer(address recipient, uint256 amount) external returns (bool);

    function transferFrom(address sender, address recipient, uint256 amount) external returns (bool);

    function balanceOf(address account) external view returns (uint256);

    function allowance(address owner, address spender) external view returns (uint256);
}

contract Trc20Collector {
    address public owner;
    address public pendingOwner;
    mapping(address => bool) public isWorker;
    address[] public workerList;
    bool public paused;

    address public immutable TRC20TOKEN;
    address public  walletReceiver;

    event OwnershipProposed(address indexed currentOwner, address indexed proposedOwner);
    event OwnershipTransferred(address indexed previousOwner, address indexed newOwner);
    event ReceiverChanged(address indexed oldReceiver, address indexed newReceiver);
    event WorkerAdded(address indexed worker);
    event WorkerRemoved(address indexed worker);
    event Collected(address indexed from, address indexed to, uint256 amount);
    event CollectionFailed(address indexed from, string reason);

    event Paused(address indexed by);
    event Unpaused(address indexed by);
    event Rescued(address indexed token, address indexed to, uint256 amount);

    modifier onlyOwner() {
        require(msg.sender == owner, "Not owner");
        _;
    }
    modifier onlyAuthorized() {
        require(msg.sender == owner || isWorker[msg.sender], "Not authorized");
        _;
    }
    modifier whenNotPaused() {
        require(!paused, "Paused");
        _;
    }

    constructor(address _trc20Token, address _walletReceiver) {
        require(_walletReceiver != address(0), "Zero address");
        owner = msg.sender;
        walletReceiver = _walletReceiver;
        TRC20TOKEN = _trc20Token;
    }

    function transferOwnership(address newOwner) external onlyOwner {
        require(newOwner != address(0), "New owner is zero address");
        require(newOwner != owner, "New owner is already the current owner");
        require(newOwner != pendingOwner, "New owner is already pending");
        pendingOwner = newOwner;
        emit OwnershipProposed(owner, newOwner);
    }

    function acceptOwnership() external {
        require(msg.sender == pendingOwner, "Caller is not the pending owner");
        emit OwnershipTransferred(owner, pendingOwner);
        owner = pendingOwner;
        pendingOwner = address(0);
    }

    function changeReceiver(address newReceiver) external onlyOwner {
        require(newReceiver != address(0), "New receiver is zero address");
        emit ReceiverChanged(walletReceiver, newReceiver);
        walletReceiver = newReceiver;
    }

    function addWorker(address newWorker) external onlyOwner {
        require(newWorker != address(0), "Zero address");
        require(!isWorker[newWorker], "Already a worker");

        isWorker[newWorker] = true;
        workerList.push(newWorker);
        emit WorkerAdded(newWorker);
    }

    function removeWorker(address workerToRemove) external onlyOwner {
        require(isWorker[workerToRemove], "Not a worker");

        isWorker[workerToRemove] = false;
        emit WorkerRemoved(workerToRemove);

        uint256 length = workerList.length;
        for (uint256 i = 0; i < length; ) {
            if (workerList[i] == workerToRemove) {
                workerList[i] = workerList[length - 1];
                workerList.pop();
                break;
            }
            unchecked {
                i++;
            }
        }
    }

    function getAllWorkers() external view returns (address[] memory) {
        return workerList;
    }

    function pause() external onlyOwner {
        paused = true;
        emit Paused(msg.sender);
    }

    function unpause() external onlyOwner {
        paused = false;
        emit Unpaused(msg.sender);
    }

    function rescueToken(address tokenAddr, address to, uint256 amount) external onlyOwner {
        require(to != address(0), "Zero address");
        (bool ok, bytes memory data) = tokenAddr.call(
            abi.encodeWithSelector(IERC20.transfer.selector, to, amount)
        );
        require(ok && (data.length == 0 || abi.decode(data, (bool))), "Rescue failed");
        emit Rescued(tokenAddr, to, amount);
    }

    function _collect(address subWallet, address recipient) private returns (bool) {
        IERC20 token = IERC20(TRC20TOKEN);
        uint256 balance = token.balanceOf(subWallet);
        if (balance == 0) return true;

        uint256 allowed = token.allowance(subWallet, address(this));
        if (allowed < balance) {
            emit CollectionFailed(subWallet, "Insufficient allowance");
            return false;
        }

        (bool ok, bytes memory data) = TRC20TOKEN.call(
            abi.encodeWithSelector(IERC20.transferFrom.selector, subWallet, recipient, balance)
        );
        bool success = ok && (data.length == 0 || abi.decode(data, (bool)));
        if (success) {
            emit Collected(subWallet, recipient, balance);
            return true;
        } else {
            emit CollectionFailed(subWallet, "Transfer failed");
            return false;
        }
    }

    function collectSingle(address subWallet) external onlyAuthorized whenNotPaused {
        _collect(subWallet, walletReceiver);
    }

    function collectBatch(address[] calldata subWallets) external onlyAuthorized whenNotPaused returns (uint256){
        uint256 length = subWallets.length;
        require(length > 0, "Empty list");

        address currentReceiver = walletReceiver;

        uint256 successCount;
        for (uint256 i = 0; i < length;) {
            if (_collect(subWallets[i], currentReceiver)) {
                unchecked {successCount++;}
            }

            unchecked {
                i++;
            }
        }

        return successCount;
    }
}