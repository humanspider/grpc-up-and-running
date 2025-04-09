package ecommerce;

import java.util.HashMap;
import java.util.Map;
import java.util.UUID;

import ecommerce.ProductInfoGrpc.ProductInfoImplBase;
import io.grpc.Status;
import io.grpc.StatusException;

public class ProductInfoImpl extends ProductInfoImplBase {
    private Map<String, ProductInfoOuterClass.Product> productMap = new HashMap<String, ProductInfoOuterClass.Product>();

    @Override
    public void addProduct(
        ProductInfoOuterClass.Product request,
        io.grpc.stub.StreamObserver<ProductInfoOuterClass.ProductID> responseObserver
    )
    { 
        UUID uuid = UUID.randomUUID();
        String randomUUIDString = uuid.toString();
        productMap.put(randomUUIDString, request);
        ProductInfoOuterClass.ProductID id = ProductInfoOuterClass.ProductID.newBuilder()
        .setValue(randomUUIDString).build();
        responseObserver.onNext(id);
        responseObserver.onCompleted();
    }

    @Override
    public void getProduct(
        ProductInfoOuterClass.ProductID request,
        io.grpc.stub.StreamObserver<ProductInfoOuterClass.Product> responseObserver
    ) {
        String id = request.getValue();
        if (productMap.containsKey(id)) {
            responseObserver.onNext(productMap.get(id));
            responseObserver.onCompleted();
        } else {
            responseObserver.onError(new StatusException(Status.NOT_FOUND));
        }
    }
}