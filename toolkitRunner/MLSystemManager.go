package main

import (
	"flag"
	"fmt"
	"time"
	"strconv"
	"github.com/TimWhiting/toolkitGo/toolkit"
	"github.com/TimWhiting/toolkitGo/learners"
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
	if len(extra) > 0{
		args.EvalExtra = extra[0];
	}
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

func (ml MLSystemManager)GetLearner(mod string, rand toolkit.Random)(toolkit.Learner, error){
	if mod == "baseline" {
		return &learners.BaselineLearner{}, nil;
	} else if mod ==("perceptron") {
		return nil,nil;
		//return perceptronLab.NewPerceptron(rand), nil;
		// else if (model.equals("neuralnet")) return new NeuralNet(rand);
		// else if (model.equals("decisiontree")) return new DecisionTree();
		// else if (model.equals("knn")) return new InstanceBasedLearner();
	}else{
		return nil,nil;
	}
}
func (ml MLSystemManager)Run(args Args)(error){
	learner, _ := ml.GetLearner(args.Learner,toolkit.Random{});
	data := toolkit.NewEmptyMatrix();
	data.LoadArff(args.Arff);
	data.Print();
	if args.Normalize {
		fmt.Println("Using normalized data\n")
		data.Normalize()
	}
	rand := toolkit.Random{};
	fmt.Println();
	fmt.Println("Dataset name:" ,args.Arff);
	fmt.Println("Number of instances: ",data.Rows());
	fmt.Println("Number of attributes: ",data.Cols());
	fmt.Println("Learning algorithm: ",args.Learner);
	fmt.Println("Evaluation method: ",args.Evaluation);
	fmt.Println();
	if args.Evaluation == "training"{
		fmt.Println("Calculating accuracy on training set...");
		features := toolkit.NewMatrix(data, 0, 0, data.Rows(), data.Cols() - 1);
		labels := toolkit.NewMatrix(data, 0, data.Cols() - 1, data.Rows(), 1);
		fmt.Println("Finished making matrices")
		confusion := toolkit.NewEmptyMatrix();
		startTime :=time.Now();
		learner.Train(features, labels);
		elapsedTime := time.Now().Sub(startTime);
		fmt.Println("Time to train (in seconds): " , elapsedTime.Seconds());
		accuracy, _ := learner.MeasureAccuracy(features, labels, confusion);
		fmt.Println("Training set accuracy: " , accuracy);
		if(args.Verbose) {
			fmt.Println("\nConfusion matrix: (Row=target value, Col=predicted value)");
			confusion.Print();
			fmt.Println("\n");
		}
	} else if args.Evaluation == "static"{
		testData := toolkit.NewEmptyMatrix();
		testData.LoadArff(args.EvalExtra);
		if args.Normalize {
			testData.Normalize(); // BUG! This may normalize differently from the training data. It should use the same ranges for normalization!
		}
		fmt.Println("Calculating accuracy on separate test set...");
		fmt.Println("Test set name: " , args.EvalExtra);
		fmt.Println("Number of test instances: " ,testData.Rows());
		features := toolkit.NewMatrix(data, 0, 0, data.Rows(), data.Cols() - 1);
		labels := toolkit.NewMatrix(data, 0, data.Cols() - 1, data.Rows(), 1);
		startTime :=time.Now();
		learner.Train(features, labels);
		elapsedTime := time.Now().Sub(startTime);
		fmt.Println("Time to train (in seconds): " , elapsedTime.Seconds());
		trainAccuracy,_ := learner.MeasureAccuracy(features, labels, toolkit.NewEmptyMatrix());
		fmt.Println("Training set accuracy: " , trainAccuracy);
		testFeatures := toolkit.NewMatrix(testData, 0, 0, testData.Rows(), testData.Cols() - 1);
		testLabels := toolkit.NewMatrix(testData, 0, testData.Cols() - 1, testData.Rows(), 1);
		confusion := toolkit.NewEmptyMatrix();
		testAccuracy,_ := learner.MeasureAccuracy(testFeatures, testLabels, confusion);
		fmt.Println("Test set accuracy: " , testAccuracy);
		if args.Verbose {
			fmt.Println("\nConfusion matrix: (Row=target value, Col=predicted value)");
			confusion.Print();
			fmt.Println("\n");
		}
	} else if args.Evaluation == "random"{
		fmt.Println("Calculating accuracy on a random hold-out set...");
		trainPercent,_ := strconv.ParseFloat(args.EvalExtra,64);
		if (trainPercent < 0 || trainPercent > 1) {
			panic("Percentage for random evaluation must be between 0 and 1");
		}
		fmt.Println("Percentage used for training: " , trainPercent);
		fmt.Println("Percentage used for testing: " ,(1 - trainPercent));
		data.Shuffle(rand);
		trainSize := int(trainPercent * float64(data.Rows()));
		trainFeatures := toolkit.NewMatrix(data, 0, 0, trainSize, data.Cols() - 1);
		trainLabels := toolkit.NewMatrix(data, 0, data.Cols() - 1, trainSize, 1);
		testFeatures := toolkit.NewMatrix(data, trainSize, 0, data.Rows() - trainSize, data.Cols() - 1);
		testLabels := toolkit.NewMatrix(data, trainSize, data.Cols() - 1, data.Rows() - trainSize, 1);
		startTime :=time.Now();
		learner.Train(trainFeatures, trainLabels);
		elapsedTime := time.Now().Sub(startTime);
		fmt.Println("Time to train (in seconds): " , elapsedTime.Seconds());
		trainAccuracy,_ := learner.MeasureAccuracy(trainFeatures, trainLabels, toolkit.NewEmptyMatrix());
		fmt.Println("Training set accuracy: " ,trainAccuracy);
		confusion := toolkit.NewEmptyMatrix();
		testAccuracy,_ := learner.MeasureAccuracy(testFeatures, testLabels, confusion);
		fmt.Println("Test set accuracy: " ,testAccuracy);
		if args.Verbose {
			fmt.Println("\nConfusion matrix: (Row=target value, Col=predicted value)");
			confusion.Print();
			fmt.Println("\n");
		}
	} else if (args.Evaluation == "cross"){
		fmt.Println("Calculating accuracy using cross-validation...");
		folds,_ := strconv.ParseInt(args.EvalExtra,10,64);
		if folds <= 0 {
			panic("Number of folds must be greater than 0");
		}
		fmt.Println("Number of folds: " , folds);
		reps := 1;
		sumAccuracy := 0.0;
		var elapsedTime time.Duration;
		for j := 0; j < reps; j++ {

			data.Shuffle(rand);
			for i := 0; i < int(folds); i++ {
				begin := i * data.Rows()/ int(folds);
				end :=(i + 1) * data.Rows() / int(folds);
				trainFeatures := toolkit.NewMatrix(data, 0, 0, begin, data.Cols() - 1);
				trainLabels := toolkit.NewMatrix(data, 0, data.Cols() - 1, begin, 1);
				testFeatures := toolkit.NewMatrix(data, begin, 0, end - begin, data.Cols() - 1);
				testLabels := toolkit.NewMatrix(data, begin, data.Cols() - 1, end - begin, 1);
				trainFeatures.Add(data, end, 0, data.Rows() - end);
				trainLabels.Add(data, end, data.Cols() - 1, data.Rows() - end);
				startTime := time.Now();
				learner.Train(trainFeatures, trainLabels);
				elapsedTime = time.Now().Sub(startTime) + elapsedTime;
				accuracy,_ := learner.MeasureAccuracy(testFeatures, testLabels, toolkit.NewEmptyMatrix());
				sumAccuracy += accuracy;
				fmt.Println("Rep=" , j , ", Fold=" , i , ", Accuracy=" , accuracy);
			}
		}
		timeTotal := elapsedTime.Seconds()/float64((reps * int(folds)));
		fmt.Println("Average time to train (in seconds): " , timeTotal);
		fmt.Println("Mean accuracy=" , (sumAccuracy / float64(reps * int(folds))));
	}
	return nil;
}