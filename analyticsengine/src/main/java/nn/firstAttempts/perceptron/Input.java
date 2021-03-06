package nn.firstAttempts.perceptron;

/**
 * Special kind of Perceptron, used as input for a nn:
 * value returns a literal, instead of a computed value.
 */
public class Input extends Perceptron {
    private double value;

    public Input(double value) {
        super(0.0);
        this.value = value;
    }

    @Override
    public double value() {
        return value;
    }

    public void setValue() {this.value = value;}
}
