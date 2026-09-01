#ifndef NN_H
#define NN_H

typedef struct {
    int dim;
    float *v;
} Vector;

typedef struct _Node {
    Vector* weights;
    float bias;
} Node;

typedef struct _Layer {
    int n_nodes;
    Node* nodes;
} Layer;

typedef struct {
    int n_layers;
    Layer* layers;
} Network;

Vector* sum(Vector* a, Vector* b);
float dot(Vector* a, Vector* b);
Vector* layer_out(Layer* l, Vector* v_in);
Vector* forward_pass(Network* nn, Vector* input);

Layer* init_layer(int n, int dim);

void print_vector(Vector* a);
void print_layer(Layer* l);

#endif
