package toolkitGo

import (
	"math/rand"
	"flag"
	"fmt"
)

type MLSystemManager struct{

}
type Args struct{
	Arff,Learner,Evaluation,EvalExtra string
	Verbose,Normalize bool
}

func main(){
	ml := MLSystemManager{};
	args := Args{}
	flag.StringVar(&args.Learner,"L", "foo", "learner")
	flag.StringVar(&args.Evaluation,"E", "foo", "evaluation")
	flag.StringVar(&args.Arff,"A", "foo", "arff")
	flag.BoolVar(&args.Verbose,"V", false, "verbose")
	flag.BoolVar(&args.Normalize,"N", false, "normalize")
	flag.Parse()
	extra := flag.Args();
	args.EvalExtra = extra[0];
	if args.Arff == "" || args.Learner == "" || args.Evaluation == "" || !(args.Evaluation == "static" || args.Evaluation == "random" || args.Evaluation == "cross" || args.Evaluation == "training" ){
		fmt.Println("Usage:")
		fmt.Println("MLSystemManager -L [learningAlgorithm] -A [ARFF_File] -E [evaluationMethod] {[extraParamters]} [OPTIONS]\n");
		fmt.Println("OPTIONS:");
		fmt.Println("-V Print the confusion matrix and learner accuracy on individual class values");
		fmt.Println("-N Use normalized data");
		fmt.Println();
		fmt.Println("Possible evaluation methods are:");
		fmt.Println("MLSystemManager -L [learningAlgorithm] -A [ARFF_File] -E training");
		fmt.Println("MLSystemManager -L [learningAlgorithm] -A [ARFF_File] -E static [testARFF_File]");
		fmt.Println("MLSystemManager -L [learningAlgorithm] -A [ARFF_File] -E random [%_ForTraining]");
		fmt.Println("MLSystemManager -L [learningAlgorithm] -A [ARFF_File] -E cross [numOfFolds]\n");
	}
	ml.Run(args);
}

func (ml MLSystemManager)GetLearner(mod string, rand rand.Rand)(SupervisedLearner, error){
	if (mod == "baseline") {
	return BaselineLearner{}, nil
	}// else if (model.equals("perceptron")) return new Perceptron(rand);
	// else if (model.equals("neuralnet")) return new NeuralNet(rand);
	// else if (model.equals("decisiontree")) return new DecisionTree();
	// else if (model.equals("knn")) return new InstanceBasedLearner();
	else{
		return nil,nil;
	}
}
func (ml MLSystemManager)Run(args Args)(error){
	return nil;
}