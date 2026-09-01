#include <stdio.h>
#include <stdlib.h>

#include "nn.h"


int main() {

    Vector* w1 = (Vector*)malloc(sizeof(Vector));
    w1->dim = 3;
    w1->v = (float*)malloc(sizeof(int)*w1->dim);
    w1->v[0] = 1;
    w1->v[1] = 2;
    w1->v[2] = 3;

    Vector* w2 = (Vector*)malloc(sizeof(Vector));
    w2->dim = 3;
    w2->v = (float*)malloc(sizeof(int)*w2->dim);
    w2->v[0] = 3;
    w2->v[1] = 4;
    w2->v[2] = 1;

    print_vector(w1);
    print_vector(w2);

    Vector* w3 = sum(w1, w2);
    int r = dot(w1, w2);

    printf("------\n");
    print_vector(w3);
    printf("dot: %d\n", r);

    Layer* l1 = init_layer(10, 3);
    Layer* l2 = init_layer(32, 10);
    Layer* l3 = init_layer(4, 32);
    print_layer(l1);
    print_layer(l2);
    print_layer(l3);

    Vector* w4 = layer_out(l1, w1);
    printf("layer out: ");
    print_vector(w4);

    Network* nn = (Network*)malloc(sizeof(Network));
    nn->n_layers = 3;
    nn->layers = (Layer*)malloc(sizeof(Layer) * nn->n_layers);
    nn->layers[0].n_nodes = l1->n_nodes;
    nn->layers[0].nodes = l1->nodes;
    nn->layers[1].n_nodes = l2->n_nodes;
    nn->layers[1].nodes = l2->nodes;
    nn->layers[2].n_nodes = l3->n_nodes;
    nn->layers[2].nodes = l3->nodes;


    Vector* w7 = forward_pass(nn, w1);
    printf("forward pass: ");
    print_vector(w7);









    return 0;
}
