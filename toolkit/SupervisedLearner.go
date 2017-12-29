package toolkit

import (
	"fmt"
	"math"
)
type Learner interface {
	Train(features, labels Matrix)error
	Predict(features, labels []float64)error
	MeasureAccuracy(features,labels,confusion Matrix)(float64,error)
}
type SupervisedLearner struct{
	Learner
}



// Before you call this method, you need to divide your data
// into a feature matrix and a label matrix.
func (s SupervisedLearner) Train(features, labels Matrix)error{
	return nil;
}
// A feature vector goes in. A label vector comes out. (Some supervised
// learning algorithms only support one-dimensional label vectors. Some
// support multi-dimensional label vectors.)
func (s SupervisedLearner) Predict(features, labels []float64)error{
	return nil;
}

// The model must be trained before you call this method. If the label is nominal,
// it returns the predictive accuracy. If the label is continuous, it returns
// the root mean squared error (RMSE). If confusion is non-NULL, and the
// output label is nominal, then confusion will hold stats for a confusion matrix.
func (s SupervisedLearner)MeasureAccuracy(features,labels,confusion Matrix)(float64,error) {
	if (features.Rows() != labels.Rows()) {
		return 0, fmt.Errorf("Expected the features and labels to have the same number of rows");
	}
	if (labels.Cols() != 1) {
		return 0, fmt.Errorf("Sorry, this method currently only supports one-dimensional labels");
	}
	if (features.Rows() == 0) {
		return 0, fmt.Errorf("Expected at least one row");
	}
	labelValues := labels.ValueCount(0);
	if (labelValues == 0) { // If the label is continuous...
		// The label is continuous, so measure root mean squared error
		pred := make([]float64, 1);
		sse := 0.0;
		for i := 0; i < features.Rows(); i++ {
			feat := features.Row(i);
			targ := labels.Row(i);
			pred[0] = 0.0; // make sure the prediction is not biassed by a previous prediction
			s.Predict(feat, pred);
			delta := targ[0] - pred[0];
			sse += (delta * delta);
		}
		return math.Sqrt(sse / float64(features.Rows())), nil;
	} else {
		// The label is nominal, so measure predictive accuracy

		confusion.SetSize(labelValues, labelValues);
		for i := 0; i < labelValues; i++ {
			confusion.SetAttrName(i, labels.AttrValue(0, i));
		}
		correctCount := 0;
		prediction := make([]float64,1);
		for i := 0; i < features.Rows(); i++ {
			feat := features.Row(i);
			targ := int(labels.Get(i, 0));
			if (targ >= labelValues) {
				return 0, fmt.Errorf("The label is out of range");
			}
			s.Predict(feat, prediction);
			pred := int(prediction[0]);
			confusion.Set(targ, pred, confusion.Get(targ, pred) + 1);
			if (pred == targ) {
				correctCount++;
			}
		}
		return float64(correctCount) / float64(features.Rows()), nil;
	}
}