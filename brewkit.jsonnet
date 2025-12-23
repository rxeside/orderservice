local project = import 'brewkit/project.libsonnet';

local appIDs = [
    'orderservice',
];

local proto = [
    'api/server/orderinternal/orderinternal.proto',
    'api/clients/productinternal/productinternal.proto',
    'api/clients/userinternal/userinternal.proto',
    'api/clients/paymentinternal/paymentinternal.proto',
];

project.project(appIDs, proto)