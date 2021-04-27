import BigNumber from "bignumber.js";
import { IProvider } from "./interface";
import { GasPrice, GasLimit, TransactionVersion, ChainID, GasPriceModifier } from "./networkParams";

/**
 * An object holding Network configuration parameters.
 */
export class NetworkConfig {
    private static default: NetworkConfig;

    /**
     * The chain ID. E.g. "1" for the Mainnet.
     */
    public ChainID: ChainID;

    /**
     * The gas required by the Network to process a byte of the {@link TransactionPayload}.
     */
    public GasPerDataByte: number;
    /**
     * The round duration.
     */
    public RoundDuration: number;
    /**
     * The number of rounds per epoch.
     */
    public RoundsPerEpoch: number;

    /**
     * The Top Up Factor for APR calculation
     */
    public TopUpFactor: number;

    /**
     * The Top Up Factor for APR calculation
     */
    public TopUpRewardsGradientPoint: BigNumber;

    /**
     *
     */
    public GasPriceModifier: GasPriceModifier;

    /**
     * The minimum gas limit required to be set when broadcasting a {@link Transaction}.
     */
    public MinGasLimit: GasLimit;

    /**
     * The minimum gas price required to be set when broadcasting a {@link Transaction}.
     */
    public MinGasPrice: GasPrice;

    /**
     * The oldest {@link TransactionVersion} accepted by the Network.
     */
    public MinTransactionVersion: TransactionVersion;

    /**
     * True if sync(provider) has been called.
     */
    public isSynchronized: boolean;

    constructor() {
        this.ChainID = new ChainID("T");
        this.GasPerDataByte = 1500;
        this.TopUpFactor = 0;
        this.RoundDuration = 0;
        this.RoundsPerEpoch = 0;
        this.TopUpRewardsGradientPoint = new BigNumber(0);
        this.MinGasLimit = new GasLimit(50000);
        this.MinGasPrice = new GasPrice(1000000000);
        this.GasPriceModifier = new GasPriceModifier(1);
        this.MinTransactionVersion = new TransactionVersion(1);
        this.isSynchronized = false;
    }

    /**
     * Gets the default configuration object (think of the Singleton pattern).
     */
    static getDefault(): NetworkConfig {
        if (!NetworkConfig.default) {
            NetworkConfig.default = new NetworkConfig();
        }

        return NetworkConfig.default;
    }

    /**
     * Synchronizes a configuration object by querying the Network, through a {@link IProvider}.
     * @param provider The provider to use
     */
    async sync(provider: IProvider): Promise<void> {
        let fresh: NetworkConfig = await provider.getNetworkConfig();
        Object.assign(this, fresh);
        this.isSynchronized = true;
    }

    /**
     * Constructs a configuration object from a HTTP response (as returned by the provider).
     */
    static fromHttpResponse(payload: any): NetworkConfig {
        let networkConfig = new NetworkConfig();

        networkConfig.ChainID = new ChainID(payload["erd_chain_id"]);
        networkConfig.GasPerDataByte = Number(payload["erd_gas_per_data_byte"]);
        networkConfig.TopUpFactor = Number(payload["erd_top_up_factor"]);
        networkConfig.RoundDuration = Number(payload["erd_round_duration"]);
        networkConfig.RoundsPerEpoch = Number(payload["erd_rounds_per_epoch"]);
        networkConfig.TopUpRewardsGradientPoint = new BigNumber(payload["erd_rewards_top_up_gradient_point"]);
        networkConfig.MinGasLimit = new GasLimit(payload["erd_min_gas_limit"]);
        networkConfig.MinGasPrice = new GasPrice(payload["erd_min_gas_price"]);
        networkConfig.MinTransactionVersion = new TransactionVersion(payload["erd_min_transaction_version"]);
        networkConfig.GasPriceModifier = new GasPriceModifier(payload["erd_gas_price_modifier"]);

        return networkConfig;
    }
}
