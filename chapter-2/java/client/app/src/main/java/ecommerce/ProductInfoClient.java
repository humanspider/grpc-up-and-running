package ecommerce;

import io.grpc.ManagedChannel;
import io.grpc.ManagedChannelBuilder;

public class ProductInfoClient {
    public static void main(String[] args) throws InterruptedException {
        ManagedChannel channel = ManagedChannelBuilder
            .forAddress("localhost", 50051)
            .usePlaintext()
            .build();

        ProductInfoGrpc.ProductInfoBlockingStub stub = ProductInfoGrpc.newBlockingStub(channel);

        String description = "Meet Apple iPhone 11. " + "All-new dual-camera system with " +
        "Ultra Wide and Night mokde.";
        ProductInfoOuterClass.ProductID productID = stub.addProduct(ProductInfoOuterClass.Product.newBuilder()
            .setName("Apple iPhone 11")
            .setDescription(description)
            .build());
        System.out.println(productID.getValue());

        ProductInfoOuterClass.Product product = stub.getProduct(productID);
        System.out.println(product.toString());
        channel.shutdown();
    }
}
