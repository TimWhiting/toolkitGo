package learners

import "github.com/TimWhiting/toolkitGo/toolkit"

type BaselineLearner struct{
	*toolkit.SupervisedLearner
	m_labels []float64
}

func (bl *BaselineLearner)Train(features, labels *toolkit.Matrix)(error){
	bl.m_labels = make([]float64,labels.Cols())
	for i := 0; i < labels.Cols(); i++{
		if labels.ValueCount(i) == 0{
			bl.m_labels[i] = labels.ColumnMean(i) // continuous
		}else{
			bl.m_labels[i] = labels.MostCommonValue(i); //nominal
		}
	}
	return nil;
}

func (bl *BaselineLearner)Predict(features, labels []float64)(error){
	for i := 0; i < len(bl.m_labels); i++{
		labels[i] = bl.m_labels[i];
	}
	return nil;
}