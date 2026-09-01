#ifndef NN_H
#define NN_H

#define EULER_NUMBER 2.71828182846

typedef struct {
    int cols;
    int rows;
    double** v;
} Matrix;

typedef struct _Layer {
    char a_func;
    Matrix* biases;
    Matrix* weights;
} Layer;

typedef struct {
    int n_layers;
    Layer* layers;
} Network;

Matrix* sum(Matrix* a, Matrix* b);
Matrix* product(Matrix* a, Matrix* b);
Matrix* transpose(Matrix* m);
Matrix** forward_pass(Network* nn, Matrix* input);

Layer* init_layer(int in_dim, int n_nodes, char a_func);

void free_matrix(Matrix* m);

void print_matrix(Matrix* a);
void print_layer(Layer* l);
void print_network(Network* nn);

#endif
