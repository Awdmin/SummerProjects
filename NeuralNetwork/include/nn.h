#ifndef NN_H
#define NN_H

typedef struct {
    int dim;
    float *v;
} Vector;

typedef struct {
    int cols;
    int rows;
    float** v;
} Matrix;

typedef struct _Node {
    Vector* weights;
    float bias;
} Node;

typedef struct _Layer_v {
    int n_nodes;
    Node* nodes;
} Layer_v;

typedef struct _Layer {
    Matrix* biases;
    Matrix* weights;
} Layer;

typedef struct {
    int n_layers;
    Layer_v* layers;
} Network_v;

typedef struct {
    int n_layers;
    Layer* layers;
} Network;

Matrix* sum(Matrix* a, Matrix* b);
float dot(Vector* a, Vector* b);
Matrix* product(Matrix* a, Matrix* b);
Matrix* transpose(Matrix* m);
Vector* layer_out_v(Layer_v* l, Vector* v_in);
Matrix* forward_pass(Network* nn, Matrix* input);

Layer_v* init_layer_v(int n, int dim);
Layer* init_layer(int n, int in_dim);

void free_matrix(Matrix* m);

void print_vector(Vector* a);
void print_matrix(Matrix* a);
void print_layer(Layer* l);

#endif
