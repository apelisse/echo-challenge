# Overview

This tutorial demonstrates how you can run a simple application that reads
through a pubsub topic from within your kubernetes cluster, with minimal set-up.

# Steps

## Set-up assumptions

We will assume that you already have a cluster set-up, with service-catalog, and
that you have created a pubsub service instance.

## Building the app

You should probably skip this step and use the published image as shown below.

If you want to modify the app, you can build it with

```
docker build -t ${your_registry}/echo .
```

And push it to your project registry:

```
gcloud docker push ${your_registry}/echo
```

You'll also need to use that location in `echo.yaml`.

## Basic echo

We have a basic echo application that receives HTTP queries, and pushes the body
of those queries into the pubsub configured with the service-catalog.

```
kubectl apply -f echo.yaml
```

## Binding the pubsub with our application

XXX: Create the binding

## Pushing data through our application

The following command will connect to our application, which will push "Hello
service!" inside the pubsub topic:

```
kubectl run --rm=true --restart=Never -i -t --image=tutum/curl curl -- curl echo -d "Hello service!"
```

## Make sure it works

Once the data has been pushed to the pubsub, we can verify that it has been
received by running this command:

```
gcloud beta pubsub subscriptions pull echo --auto-ack
```

# Conclusion

Service-catalog is awesome, and let's you seamlessly connect to other GCP
services from within your Kubernetes cluster.

# Incomplete parts

I'm not sure how to create the bindings yet, and I definitely don't read the
proper values from the secret, as I don't have a running example of the
service-catalog.
